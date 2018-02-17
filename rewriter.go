package csvrewrite

import (
	"time"
	"io"
	"bufio"
	"strings"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/charmap"
	"fmt"
)

type rewriter struct {
	transactions []*Transaction
	fromDate time.Time
	inputProc InputProcessor
	outputProc OutputProcessor
	errorWriter *io.Writer
}

func NewRewriter() *rewriter {
	return &rewriter{inputProc:NewComdirectInput(PersonalPayees{}), fromDate: time.Unix(0, 0)}
}

func (r *rewriter) SetFromDate(fromDate time.Time) {
	r.fromDate = fromDate
}

func (r *rewriter) SetInputProcessor(inputProc InputProcessor) {
	r.inputProc = inputProc
}

func (r *rewriter) SetOutProcessor(outputProc OutputProcessor) {
	r.outputProc = outputProc
}

func (r *rewriter) SetErrorWriter(errorWriter *io.Writer) {
	r.errorWriter = errorWriter
}

func (r *rewriter) ImportTransactions(input io.Reader) {
	var buf string
	var err error = nil
	for ; err == nil; {
		var p []byte
		_, err = input.Read(p)
		buf += string(p)
	}

	buf = r.inputProc.PreFilter(buf)

	inputReader := strings.NewReader(buf)
	dec := transform.NewReader(inputReader, charmap.ISO8859_15.NewDecoder())
	scanner := bufio.NewScanner(dec)
	for scanner.Scan() {
		t, err := r.inputProc.ProcessLine(scanner.Text())
		if err != nil {
			r.WriteError(err.Error())
			continue
		}
		if t != nil {
			r.transactions = append(r.transactions, t)
		}
	}

}

func (r *rewriter) ExportTransactions(output io.Writer) {
	writer := bufio.NewWriter(output)
	defer writer.Flush()

	r.outputProc.WriteHeader(writer)
	for _, t := range r.transactions {
		if t.Date.Before(r.fromDate) {
			continue
		}
		r.outputProc.Process(writer, t)
	}
}

func (r *rewriter) WriteError(s string) {
	if r.errorWriter != nil {
		fmt.Fprintln(*r.errorWriter, s)
	}
}






