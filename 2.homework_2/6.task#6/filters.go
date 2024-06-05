package main

import "fmt"

type Number interface {
	~int | ~float64 | ~float32 | ~int64 | ~int32 | ~int16 | ~int8 | ~uint | ~uint64 | ~uint32 | ~uint16 | ~uint8
}

func Filter[T any](subject string, subjectsGrade map[string]map[int][]T) map[int][]T {
	filtered, ok := subjectsGrade[subject]
	if !ok {
		fmt.Printf("Subject %s not found\n", subject)
		return make(map[int][]T)
	}
	return filtered
}

func Mean[IT, OT Number](grades []IT) OT {
	var sum IT
	for _, grade := range grades {
		sum += grade
	}
	return OT(float64(sum) / float64(len(grades)))
}
