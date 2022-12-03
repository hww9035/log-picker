package main

import (
	"fmt"
	"math"
	"time"
)

// 二进制数含有具体10字符数目
func bin101() {
	t := time.Now().UnixMilli()
	sum := 0
	for i := 1; i <= 1000000000; i++ {
		flag := 0
		for j := i; j > 0; j >>= 1 {
			if j&7 == 5 {
				flag = 1
				break
			}
		}
		if flag == 1 {
			sum++
		}
	}
	fmt.Print(sum, time.Now().UnixMilli()-t)
}

// 盘子放苹果问题
func app(x, y int) int {
	if x < 0 || y <= 0 {
		return 0
	}
	if x == 0 || x == 1 || y == 1 {
		return 1
	}
	return app(x-y, y) + app(x, y-1)
}

// 最小路径和
func minSum(arr [3][3]int) int {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if x == 0 {
				arr[x][y] = arr[x][y] + arr[x][y-1]
				continue
			}
			if y == 0 {
				arr[x][y] = arr[x][y] + arr[x-1][y]
				continue
			}
			if x != 0 && y != 0 {
				arr[x][y] = arr[x][y] + int(math.Min(float64(arr[x-1][y]), float64(arr[x][y-1])))
			}
		}
	}
	return arr[2][2]
	// fmt.Println(minSum([3][3]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}))
}

// 左上角到右下角路劲和
func waySum(arr [3][3]int) int {
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			if x == 0 && y == 0 {
				arr[x][y] = 0
				continue
			}
			if x == 0 {
				arr[x][y] = 1
				continue
			}
			if y == 0 {
				arr[x][y] = 1
				continue
			}
			if x != 0 && y != 0 {
				arr[x][y] = arr[x-1][y] + arr[x][y-1]
			}
		}
	}
	return arr[2][2]
	// fmt.Println(waySum([3][3]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}))
}

// 岛屿数目问题
func numIslands(arr [][]int) int {
	rows := len(arr)
	if rows == 0 {
		return 0
	}
	cols := len(arr[0])
	if cols == 0 {
		return 0
	}

	count := 0
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if arr[x][y] == 1 {
				count++
				dfs(x, y, arr)
			}
		}
	}

	return count
	// fmt.Println(numIslands([][]int{{0, 1, 0, 1, 1}, {1, 1, 1, 0, 0}, {1, 1, 0, 0, 1}, {0, 1, 0, 1, 1}}))
}

// 深度优先搜索
func dfs(x, y int, arr [][]int) {
	rows := len(arr)
	cols := len(arr[0])
	if x < 0 || y < 0 || x >= rows || y >= cols || arr[x][y] != 1 {
		return
	}
	arr[x][y] = 0
	dfs(x, y+1, arr)
	dfs(x, y-1, arr)
	dfs(x+1, y, arr)
	dfs(x-1, y, arr)
}

// 打家劫舍问题
func djjs(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	dp := make([]int, len(arr))
	dp[0] = arr[0]
	dp[1] = arr[0]
	if arr[0] < arr[1] {
		dp[1] = arr[1]
	}

	for i := 2; i < len(arr); i++ {
		dp[i] = int(math.Max(float64(dp[i-1]), float64(arr[i]+dp[i-2])))
	}
	return dp[len(dp)-1]
	// fmt.Print(djjs([]int{2, 7, 9, 3, 1}))
}

// 爬楼梯问题，要么一步步上去，要么跨越两步上
func stairs(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	return stairs(n-2) + stairs(n-1)
}

func main() {
}
