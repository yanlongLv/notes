package main

import (
	"context"
	"time"
)

//AsyncCall ..
func AsyncCall() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*1))
	defer cancel()
	go func(ctx context.Context) {
		time.Sleep(time.Duration(time.Second * 6))
	}(ctx)
	select {
	case <-ctx.Done():
		print("call successfully !!!!")
		return
	case <-time.After(time.Duration(time.Second * 5)):
		print("error : %s", "djdj")
		return
	}
}

func main() {
	AsyncCall()
	time.Sleep(time.Duration(time.Second * 10))
}
