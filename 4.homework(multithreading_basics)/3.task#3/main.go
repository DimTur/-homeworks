package main

import "fmt"

func goAddToChan(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()

	return out
}

// merged only 2 channels
func mergeChannels1(channel1, channel2 <-chan int) <-chan int {
	mergedChannel := make(chan int)

	go func() {
		defer close(mergedChannel)
		for {
			select {
			case val, ok := <-channel1:
				if !ok {
					channel1 = nil
					if channel2 == nil {
						return
					}
					continue
				}
				mergedChannel <- val
			case val, ok := <-channel2:
				if !ok {
					channel2 = nil
					if channel1 == nil {
						return
					}
					continue
				}
				mergedChannel <- val
			}
		}
	}()

	return mergedChannel
}

// merged only any quantity of channels. That looks better.
// func mergeChannels2(chs ...<-chan int) <-chan int {
// 	mergedChan := make(chan int)

// 	go func() {
// 		defer close(mergedChan)
// 		for _, ch := range chs {
// 			for i := range ch {
// 				mergedChan <- i
// 			}
// 		}
// 	}()

// 	return mergedChan
// }

func main() {
	nums1 := []int{1, 2, 3, 4, 5}
	nums2 := []int{6, 7, 8, 9, 10}

	ch1 := goAddToChan(nums1)
	ch2 := goAddToChan(nums2)

	mergedChannels1 := mergeChannels1(ch1, ch2)

	for val := range mergedChannels1 {
		fmt.Println(val)
	}

	// mergedChannels2 := mergeChannels2(ch1, ch2)

	// for val := range mergedChannels2 {
	// 	fmt.Println(val)
	// }
}
