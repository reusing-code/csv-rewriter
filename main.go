package main

import (
	"bufio"
	"os"

	"flag"
	"time"

	"io/ioutil"

	"strings"

	"fmt"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const dateFormat string = "02.01.2006"

type payeeSubstitution interface {
	substitute(t *Transaction)
}

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
	var inProc InputProcessor = NewComdirectInput(substitution)
	var outProc OutputProcessor = &YNABOutput{}

	inputBytes, err := ioutil.ReadFile(inputFileName)
	if err != nil {
		panic(err)
	}
	inputStr := string(inputBytes)
	inputStr = inProc.preFilter(inputStr)

	inputReader := strings.NewReader(inputStr)

	outputFile, err := os.Create(outputFileName)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	dec := transform.NewReader(inputReader, charmap.ISO8859_15.NewDecoder())
	scanner := bufio.NewScanner(dec)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	logFile, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	logWriter := bufio.NewWriter(logFile)
	defer logWriter.Flush()

	outProc.WriteHeader(writer)
	for scanner.Scan() {
		t, err := inProc.processLine(scanner.Text())
		if err != nil {
			fmt.Fprintln(logWriter, err.Error())
			continue
		}
		if t != nil {
			if t.Date.Before(fromDate) {
				continue
			}
			outProc.Process(writer, t)
		}
	}
}
