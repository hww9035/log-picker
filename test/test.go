package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"log-picker/mq/nsq"
)

func Down() {
	file, err := os.Create("/Users/huangweiwei/Downloads/go.tar.gz")
	if err != nil {
		log.Fatal(err)
		return
	}
	res, _ := http.Get("https://go.dev/dl/go1.20.3.darwin-amd64.tar.gz")

	// 方式一：buf
	buf := bufio.NewWriter(file)
	n, err := buf.ReadFrom(res.Body)
	_ = buf.Flush()
	// 方式二：copy
	// n, err := io.Copy(file, res.Body)

	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("down size:", n/1024/1024, " M")
}

func TestNsq() {
	_ = nsq.InitProducer("127.0.0.1:14150")
	go func() {
		for {
			t := time.Now().Unix()
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-1", t))
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-2", t))
			_ = nsq.PubMsg("top1", fmt.Sprint("hello-3", t))
			time.Sleep(time.Second * 3)
		}
	}()
	nsq.TestConsumer("127.0.0.1:4161", "top1", "chan1")
}

// Determinant
//
//	@Description:
//
// println(Determinant([][]int{{1}}))
// println(Determinant([][]int{{1, 3}, {2, 5}}))
// println(Determinant([][]int{{2, 5, 3}, {1, -2, -1}, {1, 3, 4}}
//
//	@param matrix
//	@return int
func Determinant(matrix [][]int) int {
	if len(matrix) <= 1 {
		return matrix[0][0]
	}
	if len(matrix) == 2 {
		return matrix[0][0]*matrix[1][1] - matrix[0][1]*matrix[1][0]
	}
	f := 0
	sum := 0
	// 递归思想
	for i := 0; i < len(matrix[0]); i++ {
		if i == 0 {
			sum = matrix[0][0] * Determinant(makeArr(&matrix, 0))
			continue
		}
		if f == 0 {
			//减法
			sum = sum - matrix[0][i]*Determinant(makeArr(&matrix, i))
			f = 1
		} else {
			//加法
			sum = sum + matrix[0][i]*Determinant(makeArr(&matrix, i))
			f = 0
		}
	}

	return sum
}

// 构造新的切片
func makeArr(src *[][]int, clo int) [][]int {
	num := len(*src)
	var data [][]int
	for i := 1; i < num; i++ {
		var tmp []int
		for j := 0; j < len((*src)[i]); j++ {
			if j != clo {
				tmp = append(tmp, (*src)[i][j])
			}
		}
		data = append(data, tmp)
	}
	return data
}

func StrRepeat(str string) {
	if len(str) == 0 {
		str = "2[ac2[b2[p]]df]10[g]hj"
	}
	sta := make([]string, 0)
	for _, v := range str {
		tmp := string(v)
		if tmp != "]" {
			sta = append(sta, tmp)
		} else {
			index, rep := doRep(sta)
			//fmt.Printf("index=%d, rep=%s\n", index, rep)
			if index >= 0 {
				sta = append(sta[:index], rep)
			}
		}
	}

	fmt.Println(strings.Join(sta, ""))
}

func doRep(ss []string) (int, string) {
	nums := ""
	reps := ""
	index := 0
	find := false
	con := true
	for i := len(ss) - 1; i >= 0; i-- {
		b := []byte(ss[i])
		if (b[0] >= 65 && b[0] <= 90) || (b[0] >= 97 && b[0] <= 122) || ss[i] == "[" {
			if find {
				con = false
			}
			continue
		}
		// 多为数字
		_, err := strconv.Atoi(ss[i])
		if err == nil && con {
			find = true
			if len(reps) == 0 {
				reps = strings.Join(ss[i+2:], "")
			}
			nums = ss[i] + nums
			index = i
		}
	}
	if len(reps) > 0 {
		num, _ := strconv.Atoi(nums)
		return index, strings.Repeat(reps, num)
	}

	return 0, ""
}
