package main

import (
	"fmt"
	"math/big"
)

func main() {
	primeChan := make(chan int)
	compChan := make(chan int)
	done := make(chan bool)

	primeNum, compNum := []int{}, []int{}

	numsArr := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 23, 29}

	go devisionNums(numsArr, primeChan, compChan, done)

	go processService(primeChan, &primeNum)
	go processService(compChan, &compNum)

	<-done

	fmt.Println("Prime numbers:", primeNum)
	fmt.Println("Composite numbers:", compNum)
}

func devisionNums(numsArr []int, primeChan, compChan chan<- int, done chan<- bool) {
	defer close(primeChan)
	defer close(compChan)
	for _, num := range numsArr {
		if big.NewInt(int64(num)).ProbablyPrime(0) {
			primeChan <- num
		} else {
			compChan <- num
		}
	}

	done <- true
}

func processService(input <-chan int, output *[]int) {
	for num := range input {
		*output = append(*output, num)
	}
}
