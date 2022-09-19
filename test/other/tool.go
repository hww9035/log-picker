package other

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	EOF = '\n'
)

func init() {
	//fmt.Println("this is tools init")
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("file close fail")
	}
}

func ReadFileString(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("ReadFileString server.txt file fail")
	}
	defer file.Close()
	read := bufio.NewReader(file)
	for {
		str, err := read.ReadString(EOF)
		if err == io.EOF {
			break
		}
		fmt.Print(str)
	}
}

func ReadFileMap(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("ReadFileBytes server.txt file fail")
	}
	defer closeFile(file)
	read := bufio.NewReader(file)
	bMap := make(map[string]int)
	for {
		str, err := read.ReadString(EOF)
		if err == io.EOF {
			break
		}
		rt := []rune(str)
		for i := 0; i < len(rt); i++ {
			s := string(rt[i])
			bMap[s] += 1
		}
	}
	fmt.Println(bMap)
}

func fbc(n int) int {
	if n < 0 {
		return 0
	}
	if n < 3 {
		return 1
	}
	return fbc(n-1) + fbc(n-2)
}

func GetRuntime() {
	res := make(map[string]string)
	res["cpu"] = strconv.Itoa(runtime.NumCPU())
	res["version"] = runtime.Version()
	res["os"] = runtime.GOOS
	res["arch"] = runtime.GOARCH
	res["compiler"] = runtime.Compiler
	fmt.Println(res)
}

// Nine 九九乘法
func Nine() {
	for i := 1; i <= 9; i++ {
		for j := 1; j <= i; j++ {
			fmt.Printf("%d * %d = %d", i, j, i*j)
			fmt.Print(" ")
		}
		fmt.Print("\n")
	}
}

func GoThread() {
	for i := 1; i <= 100; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	time.Sleep(time.Second)
	fmt.Println("GoThread end")
}

func InnerFunc() int {
	a := 1
	b := func(c int) int {
		b := c + 1
		return b
	}
	return b(a)
}

func GetOutBoundIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("localAddr:%s\n", localAddr.String())
	ip := strings.Split(localAddr.IP.String(), ":")[0]
	fmt.Println("ip:" + ip)
	return ip
}
