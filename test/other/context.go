package other

import (
	"context"
	"fmt"
	"time"
)

func cf1(ctx context.Context) {
	go cf2(ctx)
	for {
		fmt.Println("this is f1")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}

func cf2(ctx context.Context) {
	for {
		fmt.Println("this is f2")
		time.Sleep(time.Millisecond * 500)
		select {
		case <-ctx.Done():
			break
		default:
		}
	}
}

func TestCtx() {
	ctx, cancel := context.WithCancel(context.Background())
	go cf1(ctx)
	time.Sleep(time.Second * 5)
	cancel()
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

func TestContext() {
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
