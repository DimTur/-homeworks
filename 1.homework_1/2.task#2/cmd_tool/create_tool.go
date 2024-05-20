package main

import (
	"fmt"
	"os"
)

// a function that creates file with specified name
func createFile(fileName string, overwrite bool) {
	if _, err := os.Stat(fileName); err == nil {
		if !overwrite {
			fmt.Printf("File %s already exists and overwrite is not enabled.\n", fileName)
			return
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("File creating error:", err)
		return
	}
	defer file.Close()
	fmt.Printf("File '%s' successfully created.\n", fileName)
}
