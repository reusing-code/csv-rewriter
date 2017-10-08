package main

import (
	"bufio"
	"os"

	"flag"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const dateFormat string = "02.01.2006"

type payeeSubstitution interface {
	substitute(t *Transaction)
}

/**
 * Input: "Buchungstag";"Wertstellung (Valuta)";"Vorgang";"Buchungstext";"Umsatz in EUR";
 * Output: Date;Payee;Category;Memo;Outflow;Inflow
 */

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

	rewriteCSV(*inputFileName, *outputFileName, date, personalPayees{})
}

func rewriteCSV(inputFileName string, outputFileName string, fromDate time.Time, substitution payeeSubstitution) {

	inputFile, err := os.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()
	dec := transform.NewReader(inputFile, charmap.ISO8859_15.NewDecoder())
	scanner := bufio.NewScanner(dec)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	var inProc InputProcessor = NewComdirectInput(fromDate, substitution)
	var outProc OutputProcessor = &YNABOutput{}
	outProc.WriteHeader(writer)
	for scanner.Scan() {
		if t := inProc.processLine(scanner.Text()); t != nil {
			outProc.Process(writer, t)
		}
	}
}
