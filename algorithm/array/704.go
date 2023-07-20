package main

import "fmt"

func main() {
	nums := []int{0, 1, 2, 2, 3, 0, 4, 2}
	val := 2
	fmt.Println(removeElement(nums, val), nums)
}

func removeElement(nums []int, val int) int {
	length := len(nums)
	res := 0
	for i := 0; i < length; i++ {
		if nums[i] != val {
			nums[res] = nums[i]
			res++
		}
	}

	nums = nums[:res]
	return res
}
