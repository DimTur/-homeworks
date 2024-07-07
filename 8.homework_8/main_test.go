package main

import (
	"math/rand"
	"testing"
)

func BenchmarkStandardFunc(b *testing.B) {
	r := rand.New(rand.NewSource(99))
	size := 1000000
	arr := make([]int, size)
	for i := range arr {
		arr[i] = r.Intn(size)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MinEl2(arr)
	}
}

func BenchmarkConcurrensyFunc(b *testing.B) {
	r := rand.New(rand.NewSource(99))
	size := 1000000
	arr := make([]int, size)
	for i := range arr {
		arr[i] = r.Intn(size)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		MinEl3(arr)
	}
}
