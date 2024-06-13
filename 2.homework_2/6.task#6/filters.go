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

func Reduce[T1, T2 any](s []T1, init T2, f func(T1, T2) T2) T2 {
	r := init
	for _, v := range s {
		r = f(v, r)
	}
	return r
}

func Mean[IT, OT Number](grades []IT) OT {
	sum := Reduce(grades, 0, func(a IT, b IT) IT {
		return a + b
	})
	return OT(float64(sum) / float64(len(grades)))
}
