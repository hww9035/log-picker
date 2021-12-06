package utils

import (
	"fmt"
)

func init() {
	//fmt.Println("this is sort init")
}

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

// InsertSort 插入排序
func InsertSort() {

}

// QuickSort 快速排序
func QuickSort() {

}
