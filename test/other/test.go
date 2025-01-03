package other

import (
	"fmt"
	"log-picker/mq/nsq"
	"time"
	"strings"
	"strconv"
)

// Bubble 冒泡排序
func Bubble() {
	arr := [10]int{3, 5, 7, 1, 8, 2, 4, 9, 6, 10}
	for i := 1; i <= len(arr)-1; i++ {
		for j := 0; j <= len(arr)-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
	fmt.Println(arr)
}

// BinarySearch 二分查找
func BinarySearch(v int) {
	arr := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	left := 0
	right := len(arr) - 1
	mid := (right + left) / 2
	for left <= right {
		if arr[mid] == v {
			fmt.Println("找到值：", v)
			break
		} else if arr[mid] < v {
			left = mid + 1
			mid = (right + left) / 2
			continue
		} else if arr[mid] > v {
			right = mid - 1
			mid = (right + left) / 2
			continue
		}
		fmt.Println("没有找到值：", v)
	}
}

// SelSort 选择排序
func SelSort() {
	arr := [10]int{3, 5, 7, 1, 8, 2, 4, 9, 6, 10}
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	fmt.Println(arr)
}

func testNsq() {
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

// determinant
//
// println(Determinant([][]int{{1}}))
// println(Determinant([][]int{{1, 3}, {2, 5}}))
// println(Determinant([][]int{{2, 5, 3}, {1, -2, -1}, {1, 3, 4}}
//
//	@param matrix
//	@return int
func determinant(matrix [][]int) int {
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
			sum = matrix[0][0] * determinant(makeArr(&matrix, 0))
			continue
		}
		if f == 0 {
			//减法
			sum = sum - matrix[0][i]*determinant(makeArr(&matrix, i))
			f = 1
		} else {
			//加法
			sum = sum + matrix[0][i]*determinant(makeArr(&matrix, i))
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

func strRepeat(str string) {
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
			if index >= 0 {
				sta = append(sta[:index], rep)
			}
		}
	}

	fmt.Println(strings.Join(sta, ""))
}

func doRep(ss []string) (int, string) {
	numbs := ""
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
			numbs = ss[i] + numbs
			index = i
		}
	}
	if len(reps) > 0 {
		num, _ := strconv.Atoi(numbs)
		return index, strings.Repeat(reps, num)
	}

	return 0, ""
}
