package metric

//Bucket ..
type Bucket struct {
	Points []float64
	Count  int64
	next   *Bucket
}

//Window ...
type Window struct {
	window []Bucket
	size   int
}

//WindoeOpts ..
type WindoeOpts struct {
	Size int
}

//NewWindow ..
func NewWindow(opts WindoeOpts) *Window {
	buckets := make([]Bucket, opts.Size)
	for offset := range buckets {
		buckets[offset] = Bucket{Points: make([]float64, 0)}
		nextOffset := offset + 1
		if nextOffset == opts.Size {
			nextOffset = 0
		}
		buckets[offset].next = &buckets[nextOffset]
	}
	return &Window{window: buckets, size: opts.Size}
}

func (w *Window)