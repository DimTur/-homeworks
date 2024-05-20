package main

import (
	"fmt"
	"os"
)

// a function that deletes file with specified name
func deleteFile(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		fmt.Printf("File named '%s' does not exist %s\n", fileName, err)
		return
	}
	fmt.Printf("File '%s' successfully deleted.\n", fileName)
}
