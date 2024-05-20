package main

import (
	"fmt"
	"os"
)

func readFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("File named '%s' does not exist %s\n", fileName, err)
		return
	}
	fmt.Printf("File '%s' data:\n%s\n", fileName, string(data))
}
