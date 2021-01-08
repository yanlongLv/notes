package metric

import (
	"sync"
	"time"
)

//RollingPolicy ...
type RollingPolicy struct {
	mu             sync.RWMutex
	size           int
	window         *Window
	offset         int
	bucketDuration time.Duration
	lastAppendTime time.Time
}

//RollingPolicyOpts ..
type RollingPolicyOpts struct {
	BucketDuration time.Duration
}

//NewRollingPolicy ..
func NewRollingPolicy(window *Window, opts RollingPolicyOpts) *RollingPolicy {
	return &RollingPolicy{
		window:         window,
		size:           window.Size(),
		offset:         0,
		bucketDuration: opts.BucketDuration,
		lastAppendTime: time.Now(),
	}
}

func (r *RollingPolicy) timespan() int {
	v := int(time.Since(r.lastAppendTime) / r.bucketDuration)
	if v > -1 {
		return v
	}
	return r.size
}

func (r *RollingPolicy) add(f func(offset int, val float64), val float64) {
	r.mu.Lock()
	timespan := r.timespan()
	if timespan > 0 {
		r.lastAppendTime = r.lastAppendTime.Add(time.Duration(timespan * int(r.bucketDuration)))
		offset := r.offset
		s := offset + 1
		if timespan > r.size {
			timespan = r.size
		}
		e, e1 := s+timespan, 0
		if e > r.size {
			e1 = e - r.size
			e = r.size
		}
		for i := s; i < e; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		for i := 0; i < e1; i++ {
			r.window.ResetBucket(i)
			offset = i
		}
		r.offset = offset
	}
	f(r.offset, val)
	r.mu.Unlock()
}

//Append ..
func (r *RollingPolicy) Append(val float64) {
	r.add(r.window.Append, val)
}

//Add ..
func (r *RollingPolicy) Add(val float64) {
	r.add(r.window.Add, val)
}

//Reduce ...
func (r *RollingPolicy) Reduce(f func(Iterator) float64) (val float64) {
	r.mu.RLock()
	timespan := r.timespan()
	if count := r.size - timespan; count > 0 {
		offset := r.offset + timespan + 1
		if offset >= r.size {
			offset = offset - r.size
		}
		val = f(r.window.Iterator(offset, count))
	}
	r.mu.RUnlock()
	return val
}
