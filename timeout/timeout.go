package main

import (
	"context"
	"fmt"
	"time"
)

//AsyncCall ..
func AsyncCall() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	defer cancel()
	go func(ctx context.Context) {
		fmt.Print("m1")
		time.Sleep(time.Duration(time.Second * 6))
		fmt.Print("m2")
		select {
		case <-ctx.Done():
			print("call successfully !!!!")
			return
		case <-time.After(time.Duration(time.Second * 5)):
			print("error : %s", "nm")
			return
		}
	}(ctx)
}

func main() {
	AsyncCall()
	time.Sleep(time.Duration(time.Second * 10))
}
