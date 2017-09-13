package main

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const YNABDateFormat string = "01/02/2006"

type YNABOutput struct {
}

func (YNABOutput) WriteHeader(w *bufio.Writer) {
	fmt.Fprintln(w, "Date,Payee,Category,Memo,Outflow,Inflow")
}

func (YNABOutput) Process(w *bufio.Writer, t *Transaction) {
	date := t.Date.Format(YNABDateFormat)
	valueStr := strconv.FormatFloat(math.Abs(t.Value), 'f', 2, 64)
	outflow := ""
	inflow := ""
	if t.Value < 0 {
		outflow = valueStr
	} else {
		inflow = valueStr
	}
	output := strings.Join([]string{date, t.Payee, t.Category, t.Comment, outflow, inflow}, ",")
	fmt.Fprintln(w, output)
}
