package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

type InefficientLogger struct{}

func (il *InefficientLogger) Info1(msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("%s INFO: %s\n", now, msg)
	fmt.Print(logMsg)
}

type EfficientLogger struct {
	mu sync.Mutex
	w  *bufio.Writer
}

func NewEfficientLogger() *EfficientLogger {
	return &EfficientLogger{
		w: bufio.NewWriter(os.Stdout),
	}
}

func (el *EfficientLogger) Info2(msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")

	var buffer bytes.Buffer
	buffer.WriteString(now)
	buffer.WriteString(" INFO ")
	buffer.WriteString(msg)
	buffer.WriteString("\n")

	el.mu.Lock()
	el.w.WriteString(buffer.String())
	el.w.Flush() // Flush the buffer to ensure the message is written out
	el.mu.Unlock()
}

func main() {
	logger1 := &InefficientLogger{}
	for i := 0; i < 1000; i++ {
		logger1.Info1("info message")
	}

	logger2 := NewEfficientLogger()
	for i := 0; i < 1000; i++ {
		logger2.Info2("info message")
	}
}
