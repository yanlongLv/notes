package metric

import (
	"fmt"
	"time"
)

//RollingCounter ..
type RollingCounter interface {
	Metric
	Aggregation
	Timespan() int
	Reduce(func(Iterator) float64) float64
}

//RollerCounterOpts ...
type RollerCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}
type rollingCounter struct {
	policy *RollingPolicy
}

//NewRollingCounter ...
func NewRollingCounter(opts RollerCounterOpts) RollingCounter {
	window := NewWindow(WindoeOpts{Size: opts.Size})
	policy := NewRollingPolicy(window, RollingPolicyOpts{BucketDuration: opts.BucketDuration})
	return &rollingCounter{
		policy: policy,
	}
}

func (r *rollingCounter) Add(val int64) {
	if val < 0 {
		panic(fmt.Errorf("stat/matric: cannot decrease in value. val: %d", val))
	}
	r.policy.Add(float64(val))
}

func (r *rollingCounter) Reduce(f func(Iterator) float64) float64 {
	return r.policy.Reduce(f)
}

func (r *rollingCounter) Avg() float64 {
	return r.policy.Reduce(Avg)
}

func (r *rollingCounter) Min() float64 {
	return r.policy.Reduce(Min)
}

func (r *rollingCounter) Max() float64 {
	return r.policy.Reduce(Max)
}

func (r *rollingCounter) Sum() float64 {
	return r.policy.Reduce(Sum)
}

func (r *rollingCounter) Timespan() int {
	return r.policy.timespan()
}

func (r *rollingCounter) Value() int64 {
	return int64(r.Sum())
}
