package timeout

type MyQue struct {
	In  []int
	Out []int
}

func main() {
	a := MyQue{
		In: make([]int, 0),
		Out: make([]int,0),
}