// Fundamentals of Go Tasks
// Task: Sum of Numbers
// Write a Go function that takes a slice of integers as input and returns the sum of all the numbers. If the slice is empty, the function should return 0.
// [Optional]: Write a test for your function

package main

import (
	"fmt"
)

func sum_of_numbers(nums[] int) int{
	if len(nums) == 0 {
		return 0
	}
	s := 0
	for _, num := range nums {
		s += num
	}
	return s
}

func main() {
	var n, num int
	fmt.Print("How many numbers? ")
	fmt.Scan(&n)
	nums := make([]int, n)
	for i:=0; i<n; i++ {
		fmt.Printf("Enter number %d", i+1)
		fmt.Scan(&num)
		nums = append(nums, num)
	}
	fmt.Println(sum_of_numbers(nums))
}
