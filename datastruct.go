package csvrewrite

import (
	"bufio"
	"time"
)

type PayeeSubstitution interface {
	substitute(t *Transaction)
}

type Transaction struct {
	Date      time.Time
	Payee     string
	Category  string
	Comment   string
	ValueCent int
}

type InputProcessor interface {
	ProcessLine(line string) (*Transaction, error)
	PreFilter(input string) string
}

type OutputProcessor interface {
	WriteHeader(w *bufio.Writer)
	Process(w *bufio.Writer, t *Transaction)
	BatchProcess(w *bufio.Writer, t []*Transaction)
}
