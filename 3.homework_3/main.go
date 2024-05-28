package main

import "fmt"

type Formatter interface {
	Format(text string) string
}

// Plain text formatter
type plainText struct{}

func (p *plainText) Format(text string) string {
	return text
}

// Bold text formatter
type bold struct{}

func (p *bold) Format(text string) string {
	return fmt.Sprintf("**%s**", text)
}

// Code text formatter
type code struct{}

func (p *code) Format(text string) string {
	return fmt.Sprintf("'%s'", text)
}

// Italic text formatter
type italic struct{}

func (p *italic) Format(text string) string {
	return fmt.Sprintf("_%s_", text)
}

// ChainTextFormatter
type ChainTextFormatter struct {
	formatters []Formatter
}

func NewChainTextFormatter() *ChainTextFormatter {
	return &ChainTextFormatter{}
}

func (c *ChainTextFormatter) AddFormatter(formatter Formatter) {
	c.formatters = append(c.formatters, formatter)
}

func (c *ChainTextFormatter) Format(text string) string {
	for _, formatter := range c.formatters {
		text = formatter.Format(text)
	}
	return text
}

func main() {
	chainFormatter := NewChainTextFormatter()
	chainFormatter.AddFormatter(&plainText{})
	chainFormatter.AddFormatter(&bold{})
	chainFormatter.AddFormatter(&code{})
	chainFormatter.AddFormatter(&italic{})

	text := "Пример текста для форматирования"
	formattedText := chainFormatter.Format(text)
	fmt.Println("Форматированный текст:", formattedText)
}
