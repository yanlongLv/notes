package metric

import "time"

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
func NewRollingCounter(opts RollerCounterOpts) rollingCounter {
	window := NewWindow(wi)
}
