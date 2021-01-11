package bbr

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	// "github.com/go-kratos/kratos/pkg/ecode"
	limit "github.com/notes/crtos/pkg/ratelimit"
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
	rtStat          metric.RollingCounter
	inFlight        int64
	winBucketPerSec int64
	conf            *Config
	prevDrop        atomic.Value
	prevDropHit     int32
	rawMaxPASS      int64
	rawMinRt        int64
}

func (l *BBR) maxPASS() int64 {
	rawMaxPass := atomic.LoadInt64(&l.rawMaxPASS)
	if rawMaxPass > 0 && l.passStat.Timespan() < 1 {
		return rawMaxPass
	}
	rawMaxPass = int64(l.passStat.Reduce(func(iterator metric.Iterator) float64 {
		var result = 1.0
		for i := 0; iterator.Next() && i < l.conf.WinBucket; i++ {
			bucket := iterator.Bucket()
			count := 0.0
			for _, p := range bucket.Points {
				count += p
			}
			result = math.Max(result, count)
		}
		return result
	}))
	if rawMaxPass == 0 {
		rawMaxPass = 1
	}
	atomic.StoreInt64(&l.rawMaxPASS, rawMaxPass)
	return rawMaxPass
}

func (l *BBR) minRT() int64 {
	rawMinRt := atomic.LoadInt64(&l.rawMinRt)
	if rawMinRt > 0 && l.rtStat.Timespan() < 1 {
		return rawMinRt
	}
	rawMinRt = int64(math.Ceil(l.rtStat.Reduce(func(iterator metric.Iterator) float64 {
		var result = math.MaxFloat64
		for i := 1; iterator.Next() && i < l.conf.WinBucket; i++ {
			bucket := iterator.Bucket()
			if len(bucket.Points) == 0 {
				continue
			}
			total := 0.0
			for _, p := range bucket.Points {
				total += p
			}
			avg := total / float64(bucket.Count)
			result = math.Min(result, avg)
		}
		return result
	})))
	if rawMinRt <= 0 {
		rawMinRt = 1
	}
	atomic.StoreInt64(&l.rawMinRt, rawMinRt)
	return rawMinRt
}

func (l *BBR) maxFlight() int64 {
	return int64(math.Floor(float64(l.maxPASS()*l.minRT()*l.winBucketPerSec)/1000.0 + 0.5))
}

func (l *BBR) shouldDrop() bool {
	if l.cpu() < l.conf.CPUThreshold {
		prevDrop, _ := l.prevDrop.Load().(time.Duration)
		if prevDrop == 0 {
			return false
		}
		if time.Since(initTime)-prevDrop <= time.Second {
			if atomic.LoadInt32(&l.prevDropHit) == 0 {
				atomic.StoreInt32(&l.prevDropHit, 1)
			}
			inFlight := atomic.LoadInt64(&l.inFlight)
			return inFlight > 1 && inFlight > l.maxFlight()
		}
		l.prevDrop.Store(time.Duration(0))
		return false
	}
	inFlight := atomic.LoadInt64(&l.inFlight)
	drop := inFlight > 1 && inFlight > l.maxFlight()
	if drop {
		prevDrop, _ := l.prevDrop.Load().(time.Duration)
		if prevDrop != 0 {
			return drop
		}
		l.prevDrop.Store(time.Since(initTime))
	}
	return drop
}

//Allow ...
func (l *BBR) Allow(ctx context.Context, opts ...limit.AllowOption) (func(info limit.DoneInfo), error) {
	allowOpts := limit.DefaultAllowOpts()
	for _, opt := range opts {
		opt.Apply(&allowOpts)
	}
	if l.shouldDrop() {
		return nil, nil //ecode.LimitExceed
	}
	atomic.AddInt64(&l.inFlight, 1)
	stime := time.Since(initTime)
	return func(do limit.DoneInfo) {
		rt := int64((time.Since(initTime) - stime) / time.Microsecond)
		l.rtStat.Add(rt)
		atomic.AddInt64(&l.inFlight, -1)
		switch do.Op {
		case limit.Success:
			l.passStat.Add(1)
			return
		default:
			return
		}
	}, nil
}
