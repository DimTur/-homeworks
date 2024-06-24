package main

func IsEqualArrays[T comparable](arr1, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	countMap := make(map[T]int)

	for _, val := range arr1 {
		countMap[val]++
	}

	for _, val := range arr2 {
		if count, exists := countMap[val]; !exists || count == 0 {
			return false
		}
		countMap[val]--
	}

	for _, count := range countMap {
		if count != 0 {
			return false
		}
	}

	return true
}

func main() {
	arr1 := []int{1, 2, 3, 4, 5}
	arr2 := []int{5, 4, 3, 2, 1}
	arr3 := []int{1, 2, 2, 3, 4}
	arr4 := []int{1, 2, 3, 4}

	equal1 := IsEqualArrays(arr1, arr2)
	equal2 := IsEqualArrays(arr1, arr3)
	equal3 := IsEqualArrays(arr1, arr4)

	println(equal1) // Выведет true
	println(equal2) // Выведет false
	println(equal3) // Выведет false
}
