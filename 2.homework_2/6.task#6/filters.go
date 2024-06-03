package main

import "fmt"

type Number interface {
	~int | ~float64 | ~float32 | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

func Filter[T any](subject string, subjectsGrade map[string]map[int][]T) map[int][]T {
	filtered := make(map[int][]T)

	if grades, ok := subjectsGrade[subject]; ok {
		for grade, gradeValues := range grades {
			filtered[grade] = gradeValues
		}
	} else {
		fmt.Printf("Subject %s not found\n", subject)
	}

	return filtered
}

func Mean[T Number](grades []T) float64 {
	var sum T
	for _, grade := range grades {
		sum += grade
	}
	return float64(sum) / float64(len(grades))
}
