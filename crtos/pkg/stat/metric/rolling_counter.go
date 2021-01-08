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
	Reduce(func(Iterator))
}

//RollerCounterOpts ...
type RollerCounterOpts struct {
	Size           int
	BucketDuration time.Duration
}
type rollingCounter struct {
	policy *RollingCounter
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
	r.policy.Add
}

func 