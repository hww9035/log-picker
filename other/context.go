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
