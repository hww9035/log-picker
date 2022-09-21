package main

import (
	"context"
	"fmt"
	"time"
	"unsafe"
)

func test() {
	var hww int
	fmt.Println(unsafe.Sizeof(hww))
}

func gen(ctx context.Context) <-chan int {
	dst := make(chan int)
	n := 1
	go func() {
		go gen2(ctx)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("gen exit")
				return
			case dst <- n:
				n++
			}
		}
	}()
	return dst
}

func gen2(ctx context.Context) {
	select {
	case <-ctx.Done():
		fmt.Println("gen2 exit")
		return
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	for n := range gen(ctx) {
		fmt.Printf("main %d\n", n)
		if n == 5 {
			cancel()
			time.Sleep(time.Second * 3)
			break
		}
	}
}
