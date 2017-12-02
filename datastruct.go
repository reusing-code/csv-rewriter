package main

import (
	"bufio"
	"time"
)

type Transaction struct {
	Date      time.Time
	Payee     string
	Category  string
	Comment   string
	ValueCent int
}

type InputProcessor interface {
	processLine(line string) *Transaction
	preFilter(input string) string
}

type OutputProcessor interface {
	WriteHeader(w *bufio.Writer)
	Process(w *bufio.Writer, t *Transaction)
	BatchProcess(w *bufio.Writer, t []*Transaction)
}
