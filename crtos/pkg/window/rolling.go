package rolling

import (
	"sync"
	"time"

	"github.com/Go-000/Week06/window"
)

//RollingCounter ..
type RollingCounter interface {
	Timespan() int
	Add(int64)
	Value() int64
	Min() float64
	Max() float64
	Avg() float64
	Sum() float64
}

type rollingCounter struct {
	mu             sync.RWMutex
	size           int
	window         *window.Window
	lastUpdateTime time.Time
	timeDuration   time.Duration
	index          int
}

func NewRollingCounter(w *window.Window, size int, timeDuration time.Duration) RollingCounter {
	return &rollingCounter{
		window:         w,
		size:           size, //rolling 的时间长度
		lastUpdateTime: time.Now(),
		timeDuration:   timeDuration,
	}
}

func (r *rollingCounter) Timespan() int {
	return int(int(time.Since(r.lastUpdateTime) / r.timeDuration))
}

func (r *rollingCounter) Add(val float64) {
	r.mu.Lock()
	timespan := r.Timespan()
	r.lastUpdateTime = time.Now() //r.lastUpdateTime.Add(time.Duration(timespan * int(r.timeDuration)))
	if timespan > r.size {
		timespan = r.size
	}
	for i := 0; i < r.size; i++ {
		r.window.ResetBucket(i)
		offset = i
	}
	r.window.Add(r.index+1, val)

}
