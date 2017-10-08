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
}

type OutputProcessor interface {
	WriteHeader(w *bufio.Writer)
	Process(w *bufio.Writer, t *Transaction)
}
