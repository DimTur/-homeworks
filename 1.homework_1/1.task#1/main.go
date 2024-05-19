package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// init struct for our parse data to work like with objects
type Quize struct {
	Question string
	Answer   string
}

func main() {

	// you can choose filename like: go run main.go -file=yourfile.csv
	// first arg is name of our flag
	// second arg is standard file without flag
	// third arg is description of flag
	fileName := flag.String("file", "problems.csv", "CSV file to read")
	flag.Parse()

	questions, err := readCSV(*fileName)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for _, q := range questions {
		fmt.Println(q.Question)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		// fmt.Printf("Question: %s, Answer: %s\n", q.Question, q.Answer)
	}
}

// record our questions and answers to struct Queze
func recordQuestions(reader *csv.Reader) ([]Quize, error) {
	var questions []Quize

	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(record) < 2 {
			return nil, fmt.Errorf("incorrect record length: %v", record)
		}

		question := Quize{
			Question: record[0],
			Answer:   record[1],
		}

		questions = append(questions, question)
	}

	return questions, nil
}

// read our csv file (open, read, close or somthing else)
func readCSV(fileName string) ([]Quize, error) {
	file, csvFileError := os.Open(fileName)
	if csvFileError != nil {
		panic(csvFileError)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2

	return recordQuestions(reader)
}
