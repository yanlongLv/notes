package bbr

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/notes/crtos/pkg/stat/metric"
	cpustat "github.com/notes/crtos/pkg/stat/sys/cpu"
)

var (
	cpu         int64
	decay       = 0.95
	initTime    = time.Now()
	defaultConf = &Config{
		Window:       time.Second * 10,
		WinBucket:    100,
		CPUThreshold: 800,
	}
)

// Config ...
type Config struct {
	Enabled      bool
	Window       time.Duration
	WinBucket    int
	Rule         string
	Debug        bool
	CPUThreshold int64
}

type cpuGetter func() int64

// func int() {
// }

func cpuproc() {
	ticker := time.NewTicker(time.Microsecond * 250)
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			fmt.Errorf("rate/limit cpuproc error", err)
			go cpuproc()
		}
	}()
	for range ticker.C {
		//Stat
		stat := &cpustat.Stat{}
		cpustat.ReadStat(stat)
		prevCPU := atomic.LoadInt64(&cpu)
		curCPU := int64(float64(prevCPU)*decay + float64(stat.Usage)*(1.0-decay))
		atomic.StoreInt64(&cpu, curCPU)
	}
}

//BBR ...
type BBR struct {
	cpu             cpuGetter
	passStat        metric.RollingCounter
	inFlight        int64
	winBucketPerSec int64
	conf            *Config
	prevDrop        atomic.Value
	prevDropHit     int64
	rawMaxPASS      int64
	rawMinRt        int64
}

func (l *BBR) maxPASS() int64 {
	rawMaxPass := atomic.LoadInt64(&l.rawMaxPASS)
	if rawMaxPass > 0 && l.passStat.Timespan() < 1 {
		return rawMaxPass
	}
}
