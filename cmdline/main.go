package main

import (
	"os"

	"flag"
	"time"

	. "github.com/reusing-code/csvrewrite"
)

const dateFormat string = "02.01.2006"

func main() {

	inputFileName := flag.String("in", "", "input file (mandatory)")
	outputFileName := flag.String("out", "", "output file (mandatory)")
	fromDate := flag.String("fromDate", "01.01.1900", "only parse from this day and after (Format: DD.MM.YYY)")

	flag.Parse()

	if *inputFileName == "" {
		panic("No input file declared. Check usage with -h")
	}

	if *outputFileName == "" {
		panic("No output file declared. Check usage with -h")
	}

	date, err := time.Parse(dateFormat, *fromDate)
	if err != nil {
		panic(err)
	}

	rewriter := NewRewriter()
	rewriter.SetInputProcessor(NewComdirectInput(PersonalPayees{}))
	rewriter.SetOutProcessor(&YNABOutput{})
	rewriter.SetFromDate(date)

	logFile, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	rewriter.SetErrorWriter(logFile)

	reader, err := os.Open(*inputFileName)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	rewriter.ImportTransactions(reader)

	outputFile, err := os.Create(*outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	rewriter.ExportTransactions(outputFile)
}

