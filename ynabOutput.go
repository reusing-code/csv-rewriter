package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

const YNABDateFormat string = "02.01.2006"

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
		outflow = formatValue(-t.ValueCent)
	} else {
		inflow = formatValue(t.ValueCent)
	}
	output := strings.Join([]string{date, t.Payee, t.Category, t.Comment, outflow, inflow}, ",")
	fmt.Fprintln(w, output)
}

func (y *YNABOutput) BatchProcess(w *bufio.Writer, t []*Transaction) {
	for _, trans := range t {
		y.Process(w, trans)
	}
}

func formatValue(v int) string {
	neg := v < 0
	if neg {
		v = -v
	}
	str := strconv.FormatInt(int64(v), 10)
	for len(str) < 3 {
		str = "0" + str
	}
	eur := str[:len(str)-2]
	cent := str[len(str)-2:]
	result := eur + "." + cent
	if neg {
		result = "-" + result
	}
	return result
}
