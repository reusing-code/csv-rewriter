package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const YNABDateFormat string = "01/02/2006"

type YNABOutput struct {
}

func (*YNABOutput) WriteHeader(w *bufio.Writer) {
	fmt.Fprintln(w, "Date,Payee,Category,Memo,Outflow,Inflow")
}

func (*YNABOutput) Process(w *bufio.Writer, t *Transaction) {
	date := t.Date.Format(YNABDateFormat)
	outflow := ""
	inflow := ""
	if t.ValueCent < 0 {
		valueStr := strconv.FormatInt(int64(-t.ValueCent), 10)
		outflow = valueStr
	} else {
		valueStr := strconv.FormatInt(int64(t.ValueCent), 10)
		inflow = valueStr
	}
	output := strings.Join([]string{date, t.Payee, t.Category, t.Comment, outflow, inflow}, ",")
	fmt.Fprintln(w, output)
}

func (y *YNABOutput) BatchProcess(w *bufio.Writer, t []*Transaction) {
	for _, trans := range t {
		y.Process(w, trans)
	}
}
