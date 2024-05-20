package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//init flags set
	usageFileToolCmd := flag.NewFlagSet("help", flag.ExitOnError)
	createFileCmd := flag.NewFlagSet("create", flag.ExitOnError)
	readFileCmd := flag.NewFlagSet("read", flag.ExitOnError)
	deleteFileCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	//add flags for --help command
	listCommands := usageFileToolCmd.String("help", "", "What commadns can we use")

	//add flags for --create command
	createFileName := createFileCmd.String("filename", "", "File's name to create")
	createFileNameOverwrite := createFileCmd.Bool("overwrite", false, "Overwrite file if it already exists")

	//add flags for --read command
	readFileName := readFileCmd.String("filename", "", "File's name to read")

	//add flags for --delete command
	deleteFileName := deleteFileCmd.String("filename", "", "File's name to delete")

	//check for the number of arguments passed
	if len(os.Args) < 2 {
		fmt.Println("Invalid command.")
		fmt.Println("Use 'help' command.")
		return
	}

	//init 4 standard usecases for our CLI
	//that is --help, --create, --read, --delete
	switch os.Args[1] {
	case "help":
		usageFileToolCmd.Parse(os.Args[1:])
		fmt.Println("Usage:", *listCommands)
		fmt.Println("	If you need to create new file:")
		fmt.Println("		create --filename=<filename> : Create a new file")
		fmt.Println("		create --filename=<filename> --overwrite: Create a new file and overwrite fith the same name")
		fmt.Println("	If you need to read file:")
		fmt.Println("		read --filename=<filename> : Read a file")
		fmt.Println("	If you need to delete file:")
		fmt.Println("		delete --filename=<filename> : Delete a file")
	case "create":
		createFileCmd.Parse(os.Args[2:])
		fmt.Println("filename:", *createFileName)
		fmt.Println("overwrite:", *createFileNameOverwrite)
		if *createFileName == "" {
			fmt.Println("Please indicate the name of the file to create")
			fmt.Println("Use help command 'help'")
			createFileCmd.PrintDefaults()
			os.Exit(1)
		}
		createFile(*createFileName, *createFileNameOverwrite)
	case "read":
		readFileCmd.Parse(os.Args[2:])
		fmt.Println("filename:", *readFileName)
		if *readFileName == "" {
			fmt.Println("Please indicate the name of the file to delete")
			fmt.Println("Use help command 'help'")
			readFileCmd.PrintDefaults()
			os.Exit(1)
		}
		readFile(*readFileName)
	case "delete":
		deleteFileCmd.Parse(os.Args[2:])
		fmt.Println("filename:", *deleteFileName)
		if *deleteFileName == "" {
			fmt.Println("Please indicate the name of the file to delete")
			fmt.Println("Use help command 'help'")
			deleteFileCmd.PrintDefaults()
			os.Exit(1)
		}
		deleteFile(*deleteFileName)
	default:
		fmt.Println("Expected current command")
	}
}
