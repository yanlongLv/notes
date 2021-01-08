package metric

//Bucket ..
type Bucket struct {
	Points []float64
	Count  int64
	next   *Bucket
}

//Append ...
func (b *Bucket) Append(val float64) {
	b.Points = append(b.Points, val)
}

//Add ...
func (b *Bucket) Add(offset int, val float64) {
	b.Points[offset] += val
	b.Count++
}

//Reset ..
func (b *Bucket) Reset() {
	b.Points = b.Points[:0]
	b.Count = 0
}

//Next ..
func (b *Bucket) Next() *Bucket {
	return b.next
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

//ResetWindow ...
func (w *Window) ResetWindow() {
	for offset := range w.window {
		w.ResetBucket(offset)
	}
}

//ResetBucket ...
func (w *Window) ResetBucket(offset int) {
	w.window[offset].Reset()
}

//ResetBuckets ..
func (w *Window) ResetBuckets(offsets []int) {
	for _, offset := range offsets {
		w.ResetBucket(offset)
	}
}

//Add ..
func (w *Window) Add(offset int, val float64) {
	if w.window[offset].Count == 0 {
		w.window[offset].Append(val)
		return
	}
	w.window[offset].Add(0, val)
}

//Append ..
func (w *Window) Append(offset int, val float64) {
	w.window[offset].Append(val)
}

//Size ..
func (w *Window) Size() int {
	return w.size
}

//Iterator ..
func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count: count,
		cur:   &w.window[offset],
	}
}
