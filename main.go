package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"flag"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"math"
	"time"
)

const dateFormat string = "02.01.2006"

type payeeSubstitution interface {
	substitute(payee string, memo string) (string, string)
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

	hasReference := false
	headerFound := false
	for scanner.Scan() {
		inputTokens := strings.Split(scanner.Text(), ";")
		if !headerFound {
			headerFound, hasReference = checkHeader(inputTokens)
		}
		if !headerFound {
			continue
		}

		if len(inputTokens) < 5 {
			continue
		}
		outputTokens := make([]string, 6)

		if hasReference {
			inputTokens = append(inputTokens[:3], inputTokens[4:]...)
		}

		outputTokens[0] = removeQuotes(inputTokens[0])
		date, _ := time.Parse(dateFormat, outputTokens[0])
		if date.Before(fromDate) {
			continue
		}
		rawValue := strings.Replace(removeQuotes(inputTokens[4]), ".", "", -1)
		value, _ := strconv.ParseFloat(strings.Replace(rawValue, ",", ".", 1), 64)
		valueStr := strconv.FormatFloat(math.Abs(value), 'f', 2, 64)
		if value < 0 {
			outputTokens[4] = valueStr
		} else {
			outputTokens[5] = valueStr
		}

		if strings.Contains(inputTokens[2], "Lastschrift") {
			s := strings.Split(inputTokens[3], " Buchungstext: ")
			outputTokens[3] = removeQuotes(s[1])
			outputTokens[1] = removeQuotes(strings.Replace(s[0], "Auftraggeber: ", "", 1))
		} else if strings.Contains(inputTokens[2], "Wertpapiere") {
			outputTokens[3] = removeQuotes(strings.Replace(inputTokens[3], "Buchungstext: ", "", 1))
			outputTokens[1] = "Transfer: .comdirect Depot"
		} else if strings.Contains(inputTokens[2], "Überweisung") {
			s := strings.Split(inputTokens[3], " Buchungstext: ")
			outputTokens[3] = removeQuotes(s[1])
			outputTokens[1] = removeQuotes(strings.Replace(s[0], "Empfänger: ", "", 1))
		} else if strings.Contains(inputTokens[2], "Visa-Umsatz") {
			outputTokens[1] = removeQuotes(inputTokens[3])
		} else if strings.Contains(inputTokens[2], "Visa-Kartenabrechnung") {
			outputTokens[1] = "Transfer: .comdirect"
		} else {
			outputTokens[3] = removeQuotes(inputTokens[3])
		}

		outputTokens[1], outputTokens[3] = substitution.substitute(outputTokens[1], outputTokens[3])
		filterRef(outputTokens)

		fmt.Fprintln(writer, strings.Join(outputTokens, ","))
	}
}

func checkHeader(tokens []string) (bool, bool) {
	isHeader := strings.Contains(tokens[0], "Buchungstag")

	hasReference := false

	if isHeader {
		hasReference = strings.Contains(tokens[3], "Referenz")
	}
	return isHeader, hasReference
}

func removeQuotes(s string) string {
	return strings.Replace(s, "\"", "", -1)
}

func filterRef(tokens []string) {
	tokens[3] = strings.Split(tokens[3], " End-to-End-Ref.:")[0]
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
