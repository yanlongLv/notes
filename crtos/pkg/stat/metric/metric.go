package metric

//Metric ...
type Metric interface {
	Add(int64)
	Value() int64
}

//Aggregation ..
type Aggregation interface {
	Min() float64
	Max() float64
	Avg() float64
	Sum() float64
}
