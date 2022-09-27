package main

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
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

func testContext() {
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

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func md5ByString(str string) string {
	//方法一
	//data := []byte(str)
	//has := md5.Sum(data)
	//md5str1 := fmt.Sprintf("%x", has)

	//方法二
	m := md5.New()
	_, _ = io.WriteString(m, str)
	arr := m.Sum(nil)
	//md5str2 := hex.EncodeToString(arr)
	md5str2 := fmt.Sprintf("%x", arr)
	return md5str2
}

func main() {
}
