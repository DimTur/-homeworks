package main

import (
	"fmt"
	"sort"
)

func sliceIntersections(slices ...[]int) []int {
	if len(slices) == 0 {
		return []int{}
	}

	setInter := make(map[int]struct{})
	for _, num := range slices[0] {
		setInter[num] = struct{}{}
	}

	sliceToSet := func(s []int) map[int]struct{} {
		set := make(map[int]struct{})
		for _, num := range s {
			set[num] = struct{}{}
		}
		return set
	}

	for _, s := range slices[1:] {
		currentSet := sliceToSet(s)
		newSet := make(map[int]struct{})
		for num := range setInter {
			if _, found := currentSet[num]; found {
				newSet[num] = struct{}{}
			}
		}
		setInter = newSet
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

	fmt.Println(sliceIntersections(s1))
	fmt.Println(sliceIntersections(s1, s2))
	fmt.Println(sliceIntersections(s1, s2, s3))
}
