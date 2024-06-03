package main

import "fmt"

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64
}

type Numbers[T Number] []T

// sum method
func (n Numbers[T]) Sum() T {
	var sum T
	for _, v := range n {
		sum += v
	}
	return sum
}

// multiplication method
func (n Numbers[T]) Mult() T {
	var mult T = 1
	for _, v := range n {
		mult *= v
	}
	return mult
}

// is equal arrays method
func (n Numbers[T]) IsEqualArrays(n2 Numbers[T]) bool {
	if len(n) != len(n2) {
		return false
	}

	mapEqual := make(map[T]int)

	for _, v := range n {
		mapEqual[v]++
	}

	for _, v := range n2 {
		mapEqual[v]--
	}

	for _, c := range mapEqual {
		if c != 0 {
			return false
		}
	}

	return true
}

// find element method
func (n Numbers[T]) IsElement(e T) bool {
	for _, i := range n {
		if i == e {
			return true
		}
	}
	return false
}

// delete element from array by value
func (n *Numbers[T]) RemoveByValue(e T) {
	for i, v := range *n {
		if v == e {
			copy((*n)[:i], (*n)[i+1:])
			*n = (*n)[:len(*n)-1]
			return
		}
	}
}

// delete element from array by index
func (n *Numbers[T]) RemoveByIndex(ind int) {
	if ind < 0 || ind >= len(*n) {
		return
	}

	copy((*n)[ind:], (*n)[ind+1:])
	*n = (*n)[:len(*n)-1]
}

func main() {
	nums := Numbers[float64]{1.9, 2.2, 3.65, 3.5, 8.111}
	nums2 := Numbers[float64]{1.5, 2.5, 3.5}
	fmt.Println("Sum elements:", nums.Sum())
	fmt.Println("Mult elements:", nums.Mult())
	fmt.Println("IsEqualArrays elements:", nums.IsEqualArrays(nums2))
	fmt.Println("isElement elements:", nums.IsElement(1))

	fmt.Println("Array before remove:", nums)
	nums.RemoveByValue(8.111)
	fmt.Println("Array after removing:", nums)

	fmt.Println("Array before remove:", nums)
	nums.RemoveByIndex(1)
	fmt.Println("Array after removing:", nums)
}
