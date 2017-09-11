package main

import (
	"io"
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
	WriteHeader(w *io.Writer)
	Process(w *io.Writer, t Transaction)
}
