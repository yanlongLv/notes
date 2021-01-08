package main

import (
	"fmt"
	"time"
)

func main() {
	time1 := time.Now()
	//duration := 1
	time.Sleep(time.Second * 1)
	a := time.Since(time1)
	fmt.Print(a)
}
