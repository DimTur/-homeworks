package main

import (
	"fmt"
	"sync"
)

// Задаёт максимальную глубину рекурсии, при которой разрешено параллельное выполнение.
const MaxParallelDepth = 10

// Результаты выполнения без глубины рекурсии
/*
ubuntu@ubuntu:~/Desktop/go_to_middle/homeworks/8.homework_8$ go test -bench=. -benchmem -benchtime=5s
goos: linux
goarch: amd64
pkg: algs
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkStandardFunc-16             948           6266863 ns/op               0 B/op          0 allocs/op
BenchmarkConcurrensyFunc-16           19         299315642 ns/op        131510651 B/op   5036564 allocs/op
PASS
ok      algs    16.753s
*/
func MinEl2(a []int) int {
	// только для первичной проверки
	if len(a) == 0 {
		return 0
	}
	if len(a) == 1 {
		return a[0]
	}
	t1 := MinEl2(a[:len(a)/2])
	t2 := MinEl2(a[len(a)/2:])
	if t1 <= t2 {
		return t1
	}
	return t2
}

// В данном варианте мы вводим константу MaxParallelDepth для ограничения глубины рекурсии.
// После достижения максимальной глубины рекурсии рекурсивные вызовы становятся последовательными.
// Если этого не сделать, то мы получим слишком большие накладные расходы, т.е. создание большого количества горутин.
// Соответсвенно это увеличивает время выполнения и применение горутин становится плохим решением для улучшения.
// Результаты выполнения с глубиной рекурсии
/*
ubuntu@ubuntu:~/Desktop/go_to_middle/homeworks/8.homework_8$ go test -bench=. -benchmem -benchtime=5s
goos: linux
goarch: amd64
pkg: algs
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkStandardFunc-16             955           6256499 ns/op               0 B/op          0 allocs/op
BenchmarkConcurrensyFunc-16         4801           1232525 ns/op          163997 B/op       5118 allocs/op
PASS
ok      algs    12.695s
*/
func minEl3Helper(a []int, depth int) int {
	if len(a) == 0 {
		return 0
	}
	if len(a) == 1 {
		return a[0]
	}

	if depth >= MaxParallelDepth {
		t1 := minEl3Helper(a[:len(a)/2], depth+1)
		t2 := minEl3Helper(a[len(a)/2:], depth+1)
		if t1 <= t2 {
			return t1
		}
		return t2
	}

	var t1, t2 int
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		t1 = minEl3Helper(a[:len(a)/2], depth+1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		t2 = minEl3Helper(a[len(a)/2:], depth+1)
	}()

	wg.Wait()

	if t1 <= t2 {
		return t1
	}
	return t2
}

func MinEl3(a []int) int {
	return minEl3Helper(a, 0)
}

func main() {
	arr1 := []int{1, 3, 8678, 4, 53456, 234, 4234}

	fmt.Println(MinEl2(arr1))

	fmt.Println(MinEl3(arr1))
}
