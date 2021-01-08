package metric

import "fmt"

//Iterator ...
type Iterator struct {
	count         int
	iteratedCount int
	cur           *Bucket
}

//Next ...
func (i *Iterator) Next() bool {
	return i.count != i.iteratedCount
}

//Bucket ..
func (i *Iterator) Bucket() Bucket {
	if !(i.Next()) {
		panic(fmt.Errorf("stat/metric: iteration out of range iteratedCount :%d count: %d", i.iteratedCount, i.count))
	}
	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.Next()
	return bucket
}
