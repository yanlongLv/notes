package bbr

import (
	"time"

	"github.com/rs/zerolog/log"
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

func int() {
}

func cpuproc() {
	ticker := time.NewTicker(time.Microsecond * 250)
	defer func() {
		ticker.Stop()
		if err := recover(); err != nil {
			log.Error("rate/limit cpuproc error", err)
			go cpuproc()
		}
	}()
	for range ticker.C {
		stat := &cpu
	}
}
