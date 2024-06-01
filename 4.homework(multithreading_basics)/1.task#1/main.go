package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// you can choose the file name if you need
	fileName := "output.csv"

	dataChannel := make(chan string)
	signalChannel := make(chan os.Signal, 1)

	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	createFile(fileName)

	go reader(dataChannel)
	go writer(dataChannel, fileName)

	<-signalChannel
	fmt.Println("Program completed successfully.")
	close(dataChannel)
}

func reader(ch chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Write a line that is added to the file: ")
		if scanner.Scan() {
			input := scanner.Text()
			ch <- input
		} else {
			break
		}
	}
}

func writer(ch <-chan string, fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("File not found: ", err)
		return
	}
	defer closeFile(file)

	for input := range ch {
		_, err := file.WriteString(input + "\n")
		if err != nil {
			fmt.Println("File writing error: ", err)
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
