package other

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"sort"
	"time"
)

func ComputeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	sha := hex.EncodeToString(h.Sum(nil))
	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
}

func Md5ByString(str string) string {
	// 方法一
	// data := []byte(str)
	// has := md5.Sum(data)
	// fmt.Sprintf("%x", has)

	// 方法二
	m := md5.New()
	_, _ = io.WriteString(m, str)
	arr := m.Sum(nil)
	// hex.EncodeToString(arr)
	return fmt.Sprintf("%x", arr)
}

// 二进制数含有具体101字符数目
func Bin101() {
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

// 递归-盘子放苹果问题
func Apple(x, y int) int {
	if x < 0 || y <= 0 {
		return 0
	}
	if x == 0 || x == 1 || y == 1 {
		return 1
	}
	return Apple(x-y, y) + Apple(x, y-1)
}

// 最小路径和
func MinSum(arr [3][3]int) int {
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
func WaySum(arr [3][3]int) int {
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
func NumIslands(arr [][]int) int {
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

// 动态规划-打家劫舍问题
func Djjs(arr []int) int {
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

// 动态规划-梅花桩问题
func Mhz(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	dp := make([]int, len(arr))
	dp[0] = 1
	for i := 1; i < len(arr); i++ {
		dp[i] = 1
		for j := 0; j < i; j++ {
			if arr[j] < arr[i] {
				dp[i] = int(math.Max(float64(dp[i]), float64(dp[j]+1)))
			}
		}
	}
	return dp[len(arr)-1]
	// fmt.Print(mhz([]int{2, 5, 1, 5, 4, 5, 5, 6}))
}

// 动态规划-爬楼梯问题，要么一步步上去，要么跨越两步上
func Stairs(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	return Stairs(n-2) + Stairs(n-1)
}

// 动态规划-最大子串和
func MaxSumChild(arr []int) int {
	length := len(arr)
	if length == 0 {
		return 0
	}
	dp := make([]int, length)
	dp[0] = arr[0]
	for i := 1; i < length; i++ {
		tmp := dp[i-1] + arr[i]
		dp[i] = tmp
		if tmp < arr[i] {
			dp[i] = arr[i]
		}
	}
	sort.Ints(dp)
	return dp[length-1]
	// fmt.Print(maxSumChild([]int{-2, -1, -3, 4, -1, 2, 1, -5}))
}

// 动态规划-最大子串乘积
func MaxMultiplyChild(arr []int) int {
	if len(arr) == 0 {
		return 0
	}
	maxVal := arr[0]
	minVal := arr[0]
	res := 0

	for i := 1; i < len(arr); i++ {
		preMax := maxVal
		max := math.Max(float64(arr[i]*maxVal), float64(arr[i]*minVal))
		maxVal = int(math.Max(float64(arr[i]), max))
		min := math.Min(float64(arr[i]*preMax), float64(arr[i]*minVal))
		minVal = int(math.Min(float64(arr[i]), min))
		res = int(math.Max(float64(maxVal), float64(res)))
	}

	return res
	// fmt.Print(maxMultiplyChild([]int{-2, -1, -3, 4, -1, 2, 1, -5}))
}

// 滑动窗口找最小覆盖子串
func MinWindow(s string, t string) string {
	lenS := len(s)
	resL, resR := -1, -1
	maxLen := lenS
	wChar, tChar := map[byte]int{}, map[byte]int{}
	for i := 0; i < len(t); i++ {
		tChar[t[i]]++
	}

	// 范围检查
	check := func() bool {
		for k, v := range tChar {
			if wChar[k] < v {
				return false
			}
		}
		return true
	}

	// 有边界移动达到满足要求，扩大范围
	for l, r := 0, 0; r < lenS; r++ {
		if r < lenS && tChar[s[r]] > 0 {
			wChar[s[r]]++
		}
		// 满足要求，开始不断左边界移动，缩小范围
		for check() && l <= r {
			if _, ok := tChar[s[l]]; ok {
				wChar[s[l]] -= 1
			}
			// 设定最小区间游标
			if r-l+1 <= maxLen {
				maxLen = r - l + 1
				resL = l
				resR = r + 1
			}
			l++
		}
	}
	if resL > -1 {
		return s[resL:resR]
	}

	return ""
}

func ZeroSolution(A []int) int {
	alength := len(A)
	total := 0
	for k := 1; k <= alength; k++ {
		max := alength - k
		for i := 0; i <= max; i++ {
			if total > 100000 {
				return -1
			}
			if sum(&A, i, k) == 0 {
				total += 1
			}
		}
	}

	return total
}

func sum(s *[]int, i, num int) int {
	sum := 0
	for j := 0; j < num; j++ {
		sum += (*s)[j+i]
	}
	return sum
}

func ReverseSolution(N int) {
	var enable_print int
	enable_print = N % 10
	for N > 0 {
		if enable_print == 0 && N%10 != 0 {
			enable_print = 1
		}
		if enable_print == 1 {
			fmt.Print(N % 10)
		}
		N = N / 10
	}
}
