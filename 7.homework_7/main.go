package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"
)

type InefficientLogger struct {
}

func NewIneffectiveLogger() *InefficientLogger {
	return &InefficientLogger{}
}

func (il *InefficientLogger) Info1(msg string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	logMsg := fmt.Sprintf("%s INFO: %s", now, msg)
	fmt.Println(logMsg)
}

type EfficientLogger struct {
	w *bufio.Writer
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

	el.w.WriteString(buffer.String())
	el.w.Flush()
}
