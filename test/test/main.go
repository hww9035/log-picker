package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var hww int
	fmt.Println(unsafe.Sizeof(hww))
}
