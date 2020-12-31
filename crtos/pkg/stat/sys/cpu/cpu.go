package cpu

import (
	"sync/atomic"
	"time"
)

const (
	interval time.Duration = time.Microsecond * 500
)

var (
	stats CPU
	usage uint64
)

//CPU ..
type CPU interface {
	Usage() (u uint64, e error)
	Info() Info
}

//Stat ...
type Stat struct {
	Usage uint64
}

//Info ..
type Info struct {
	Frequency uint64
	Quota     float64
}

//ReadStat ..
func ReadStat(stat *Stat) {
	stat.Usage = atomic.LoadUint64(&usage)
}
