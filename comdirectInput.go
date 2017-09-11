package main

import (
	"strconv"
	"strings"
)

type ComdirectInput struct {
	headerFound bool
}

func (c ComdirectInput) processLine(line string) *Transaction {
	if !c.headerFound {
		if strings.EqualFold(line, `"Buchungstag";"Wertstellung (Valuta)";"Vorgang";"Buchungstext";"Umsatz in EUR";`) ||
			strings.EqualFold(line, `"Buchungstag";"Umsatztag";"Vorgang";"Referenz";"Buchungstext";"Umsatz in EUR";`) {
			c.headerFound = true
		}
		return nil
	}
	tokens := strings.Split(line, ";")
	length := len(tokens)
	if length < 5 {
		return nil
	}
	buchungsTag := removeQuotes(tokens[0])
	vorgang := removeQuotes(tokens[2])
	buchungsText := removeQuotes(tokens[3])
	umsatz := removeQuotes(tokens[4])
	if length == 6 {
		buchungsText = removeQuotes(tokens[4])
		umsatz = removeQuotes(tokens[5])
	}
	t := Transaction{}

	return &t
}

func parseValue(str string) float64 {
	rawValue := strings.Replace(removeQuotes(str), ".", "", -1)
	value, _ := strconv.ParseFloat(strings.Replace(rawValue, ",", ".", 1), 64)
	return value
}
