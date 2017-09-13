package main

import (
	"bufio"
	"time"
)

type Transaction struct {
	Date     time.Time
	Payee    string
	Category string
	Comment  string
	Value    float64
}

type InputProcessor interface {
	processLine(line string) *Transaction
}

type OutputProcessor interface {
	WriteHeader(w *bufio.Writer)
	Process(w *bufio.Writer, t *Transaction)
}
