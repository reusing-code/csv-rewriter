package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"math"
)

/**
 * Input: "Buchungstag";"Wertstellung (Valuta)";"Vorgang";"Buchungstext";"Umsatz in EUR";
 * Output: Date;Payee;Category;Memo;Outflow;Inflow
 */

func main() {

	if len(os.Args) < 3 {
		panic("Not enough command line arguments")
	}
	inputFileName := os.Args[1]
	outputFileName := os.Args[2]

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

	for scanner.Scan() {
		inputTokens := strings.Split(scanner.Text(), ";")
		if len(inputTokens) < 5 {
			continue
		}
		outputTokens := make([]string, 6)
		if strings.Contains(inputTokens[0], "Buchungstag") {
			outputTokens = strings.Split("Date;Payee;Category;Memo;Outflow;Inflow", ";")
		} else {
			outputTokens[0] = inputTokens[0]
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
			} else if strings.Contains(inputTokens[2], "Überweisung") {
				s := strings.Split(inputTokens[3], " Buchungstext: ")
				outputTokens[3] = removeQuotes(s[1])
				outputTokens[1] = removeQuotes(strings.Replace(s[0], "Empfänger: ", "", 1))
			} else {
				outputTokens[3] = removeQuotes(inputTokens[3])
			}
		}

		doPayeeFiltering(outputTokens)
		filterRef(outputTokens)

		fmt.Fprintln(writer, strings.Join(outputTokens, ";"))
	}

}

var markets = []struct {
	in  string
	out string
}{
	{"E-CENTER", "Edeka"},
	{"EDEKA", "Edeka"},
	{"ROSSMANN", "Rossmann"},
	{"ALDI", "Aldi"},
	{"ERGO", "Ergo Direkt"},
	{"COSMOS", "Cosmos Leben"},
	{"1u1", "1und1"},
	{"VOLKSWAGEN", "Volkswagen"},
	{"COM-IN", "COM-IN"},
	{"Anna Roiser", "Anna Roiser"},
	{"Norma", "Norma"},
	{"McFit", "McFit"},
	{"Stadtwerke Ingolstadt", "Stadtwerke Ingolstadt"},
	{"e.solutions", "Einkommen"},
	{"Willner", "Willner"},
	{"REWE", "REWE"},
	{"C+A", "C+A"},
	{"Schuh Muecke", "Schuh Mücke"},
	{"E-TANKEN", "Tankstelle"},
	{"Bauhaus", "Bauhaus"},
}

func removeQuotes(s string) string {
	return strings.Replace(s, "\"", "", -1)
}

func doPayeeFiltering(tokens []string) {
	for _, market := range markets {
		if CaseInsensitiveContains(tokens[1], market.in) {
			tokens[1] = market.out
			tokens[3] = ""
		}
	}
}

func filterRef(tokens []string) {
	tokens[3] = strings.Split(tokens[3], " End-to-End-Ref.:")[0]
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
