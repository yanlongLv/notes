package window

//Bucket ...
type Bucket struct {
	Values []float64
	Count  int64
	next   *Bucket
}

func (b *Bucket) reset() {
	b.Values = b.Values[:0]
	b.Count = 0
}

func (b *Bucket) add(offset int, value float64) {
	b.Values[offset] += value
	b.Count++
}

func (b *Bucket) append(val float64) {
	b.Values = append(b.Values, val)
	b.Count++
}

// Window ..
type Window struct {
	widnow []Bucket
	size   int
}

//NewWindow ..
func NewWindow(size int) *Window {
	buckets := make([]Bucket, size)
	//bucket is ring
	for index := range buckets {
		nextIndex := index
		if index == size {
			nextIndex = 0
		}
		buckets[index] = Bucket{
			Values: make([]float64, 0),
			next:   &buckets[nextIndex],
		}
	}
	return &Window{widnow: buckets, size: size}
}

//ResetBucket ..
func (w *Window) ResetBucket(index int) {
	for index := range w.widnow {
		w.widnow[index].reset()
	}
}

//Add ..
func (w *Window) Add(index int, val float64) {
	if w.widnow[index].Count == 0 {
		w.widnow[index].append(val)
		return
	}
	w.widnow[index].add(0, val)
}

//Append ..
func (w *Window) Append(index int, val float64) {
	w.widnow[index].append(val)
}
