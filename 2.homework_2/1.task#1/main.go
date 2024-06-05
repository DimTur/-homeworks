package main

import (
	"fmt"
	"sort"
)

func sliceIntersections(slices ...[]int) []int {
	if len(slices) == 0 {
		return []int{}
	}

	for _, slice := range slices {
		if len(slice) == 0 {
			return []int{}
		}
	}

	setInter := make(map[int]struct{})
	for _, num := range slices[0] {
		setInter[num] = struct{}{}
	}

	for _, slice := range slices[1:] {
		temp := make(map[int]struct{})
		for _, num := range slice {
			if _, exists := setInter[num]; exists {
				temp[num] = struct{}{}
			}
		}
		setInter = temp
	}

	result := make([]int, 0, len(setInter))
	for num := range setInter {
		result = append(result, num)
	}
	sort.Ints(result)

	return result
}

func main() {
	s1 := []int{1, 2, 3, 2}
	s2 := []int{3, 2}
	s3 := []int{}
	s4 := []int{}

	fmt.Println(sliceIntersections(s1))
	fmt.Println(sliceIntersections(s1, s2))
	fmt.Println(sliceIntersections(s1, s2, s3))
	fmt.Println(sliceIntersections(s3, s4))
}
