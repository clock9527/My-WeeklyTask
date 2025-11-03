package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// 只出现一次的数字
func singleNumber(nums []int) int {
	m := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if val, ok := m[nums[i]]; !ok {
			m[nums[i]] = nums[i]
		} else {
			delete(m, val)
		}
	}
	var num int
	for k := range m {
		num = k
	}
	return num
}

// 回文数 字符串法
func isPalindromeNumber(num int) bool {
	if num < 0 {
		return false
	}
	numStr := strconv.Itoa(num)

	rStr := ""
	for _, v := range numStr {
		rStr = string(v) + rStr
	}
	return numStr == rStr
}

// 回文数 数值法
func isPalindromeNumber2(x int) bool {
	if x < 0 {
		return false
	}
	rNum, mNum := 0, x
	for mNum > 0 {
		rNum += mNum % 10
		if mNum /= 10; mNum == 0 {
			break
		}
		rNum *= 10
	}
	return x == rNum
}

// 有效的括号
func validParentheses(s string) bool {
	if len(s) == 0 {
		return false
	}

	pMap := map[string]string{")": "(", "]": "[", "}": "{"}
	stack := make([]string, 0, len(s)/2)
	for _, p := range s {
		if v, _ := pMap[string(p)]; len(stack) == 0 || stack[len(stack)-1] != v {
			stack = append(stack, string(p))
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

// 最长公共前缀
func longestCommonPrefix(strs []string) string {

	var minLenStr = strs[0]
	for _, sEle := range strs {
		if len(minLenStr) > len(sEle) {
			minLenStr = sEle
		}
	}

	for i := 0; i < len(strs); i++ {
		for !strings.HasPrefix(strs[i], minLenStr) {
			if len(minLenStr) == 0 {
				return ""
			}
			minLenStr = minLenStr[:len(minLenStr)-1]
		}
	}
	return minLenStr
}

// 数组 加一
func plusOne(digits []int) []int {
	addV := 1
	for i := len(digits) - 1; addV == 1; i-- {
		digits[i] += addV
		if digits[i] == 10 {
			digits[i] %= 10
			if i == 0 {
				digits = append([]int{1}, digits...)
				addV = 0
			}
		} else {
			addV = 0
		}
	}
	return digits
}

// 删除有序数组中的重复项
func removeDuplicatesFromSortedArray(nums []int) (int, []int) {
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			nums[i] = nums[j]
		}
	}
	// nums = nums[:i+1]
	// fmt.Println(nums)
	return i + 1, nums[:i+1]
}

// 合并区间
func mergeIntervals(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	var midIntervals = [][]int{intervals[0]}
	var p = 0
	for _, v := range intervals[1:] {
		if midIntervals[p][1] >= v[0] && midIntervals[p][0] <= v[0] {
			if midIntervals[p][1] < v[1] {
				midIntervals[p][1] = v[1]
			}
		} else {
			midIntervals = append(midIntervals, v)
			p++
		}
	}
	return midIntervals
}

// 两数之和-for嵌套
func twoSum1(nums []int, target int) []int {
	res := []int{}
	for i, vi := range nums {
		for j, vj := range nums[i+1:] {
			if vi+vj == target {
				res = append(res, i, j+i+1)
				return res
			}

		}
	}
	return res
}

// 两数之和-for+map
func twoSum2(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, v := range nums {
		if val, ok := numMap[target-v]; ok {
			return []int{val, i}
		}
		numMap[v] = i
	}
	return []int{}
}
func main() {

	// fmt.Println("只出现一次的数字:", singleNumber([]int{4, 1, 2, 1, 2}))

	// fmt.Println("回文数-字符串法:",isPalindromeNumber(1234321))

	// fmt.Println("回文数-数字法::",isPalindromeNumber2(121))

	// fmt.Println("有效的括号:", validParentheses("()"))

	// fmt.Println("最长公共前缀:",longestCommonPrefix([]string{"flower", "flow", "flight"}))

	// fmt.Println("数组 加一:",plusOne([]int{9, 9}))

	// l, arr := removeDuplicatesFromSortedArray([]int{1, 1, 2, 3, 3})
	// fmt.Println("删除有序数组中的重复项:", l, arr)

	// fmt.Println("合并区间:", mergeIntervals([][]int{{1, 4}, {2, 3}, {8, 10}, {5, 6}, {7, 8}}))

	fmt.Println("两数之和：", twoSum1([]int{-1, -2, -3, -4, -5}, -8))
	fmt.Println("两数之和：", twoSum2([]int{-1, -2, -3, -4, -5}, -8))
}
