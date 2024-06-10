package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// you can choose the file name if you need
	fileName := "output.csv"

	ctx, cancel := context.WithCancel(context.Background())

	dataChannel := make(chan string)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	createFile(fileName)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		reader(ctx, dataChannel)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		writer(ctx, dataChannel, fileName)
	}()

	<-signalChannel
	fmt.Println("Received termination signal. Exiting...")

	cancel()

	// wg.Wait()

	fmt.Println("Program completed successfully.")
}

func reader(ctx context.Context, ch chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Write a line that is added to the file: ")
		if scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case ch <- scanner.Text():
			}
		} else {
			break
		}
	}
}

func writer(ctx context.Context, ch <-chan string, fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("File not found: ", err)
		return
	}
	defer closeFile(file)

	for {
		select {
		case input, ok := <-ch:
			if !ok {
				return
			}
			_, err := file.WriteString(input + "\n")
			if err != nil {
				fmt.Println("File writing error: ", err)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func createFile(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Println("File creating error: ", err)
			return
		}
		defer closeFile(file)
	} else if err != nil {
		fmt.Println("Error with file: ", err)
		return
	}
}

func closeFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Println("Error closing file: ", err)
	}
}
