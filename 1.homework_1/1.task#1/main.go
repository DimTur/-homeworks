package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const (
	_  = iota
	KB = 1 << (10 * iota)
	MB
	GB
)

const MaxFileSize = 2 * GB

// init struct for our parse data to work like with objects
type Quiz struct {
	Question string
	Answer   string
}

func main() {

	// you can choose filename like: go run main.go -file=yourfile.csv
	// first arg is name of our flag
	// second arg is standard file without flag
	// third arg is description of flag
	qustionsMix := flag.Bool("mix", false, "Mix the quiz questins")
	fileName := flag.String("file", "problems.csv", "CSV file to read")
	flag.Parse()

	fileInfo, err := os.Stat(*fileName)
	if err != nil {
		panic(err)
	}

	if fileInfo.Size() > MaxFileSize {
		fmt.Println("File is too large. Maximum allowed size is 2 GB.")
		return
	}

	questions, err := readCSV(*fileName)
	if err != nil {
		panic(err)
	}

	if *qustionsMix {

		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	result := workHard(questions)

	fmt.Println(result)
}

// requst all user's answer from terminal, calculate results and return it
func workHard(questions []Quiz) string {
	reader := bufio.NewReader(os.Stdin)

	correctCounter := 0
	wrongCounter := 0
	for _, q := range questions {
		fmt.Println(q.Question)
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		userInput = strings.ToLower(userInput)
		answer := strings.ToLower(q.Answer)

		if answer == userInput {
			correctCounter++
		} else {
			wrongCounter++
		}
	}

	totalQuestions := len(questions)
	winRate := float64(correctCounter) / float64(totalQuestions) * 100

	result := fmt.Sprintf(
		"Correct answers: %d\nWrong answers: %d\nWin rate: %.2f%%\n",
		correctCounter,
		wrongCounter,
		winRate,
	)

	return result
}

// record our questions and answers to struct Queze
func recordQuestions(file *os.File) ([]Quiz, error) {
	var questions []Quiz
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lastComma := strings.LastIndex(line, ",")
		if lastComma == -1 {
			return nil, fmt.Errorf("incorrect format: %s", line)
		}

		question := line[:lastComma]
		answer := line[lastComma+1:]

		// record, err := reader.Read()
		// if err != nil {
		// 	if err == io.EOF {
		// 		break
		// 	}
		// 	return nil, err
		// }

		// if len(record) < 2 {
		// 	return nil, fmt.Errorf("incorrect record length: %v", record)
		// }

		// question := Quiz{
		// 	Question: question,
		// 	Answer:   answer,
		// }

		questions = append(questions, Quiz{
			Question: question,
			Answer:   answer,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return questions, nil

	// questions = append(questions, question)
}

// read our csv file (open, read, close or somthing else if needed)
func readCSV(fileName string) ([]Quiz, error) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// reader := csv.NewReader(file)
	// reader.FieldsPerRecord = 3

	return recordQuestions(file)
}
