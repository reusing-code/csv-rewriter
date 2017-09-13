package main

import (
	"strconv"
	"strings"
	"time"
)

const comdirectDateFormat string = "02.01.2006"

type ComdirectInput struct {
	headerFound bool
	fromDate    time.Time
	sub         payeeSubstitution
}

type handler func(*Transaction) bool

func NewComdirectInput(fromDate time.Time, sub payeeSubstitution) ComdirectInput {
	com := ComdirectInput{}
	com.headerFound = false
	com.fromDate = fromDate
	com.sub = sub
	return com
}

func (c ComdirectInput) processLine(line string) *Transaction {
	if !c.headerFound {
		if strings.EqualFold(line, `"Buchungstag";"Wertstellung (Valuta)";"Vorgang";"Buchungstext";"Umsatz in EUR";`) ||
			strings.EqualFold(line, `"Buchungstag";"Umsatztag";"Vorgang";"Referenz";"Buchungstext";"Umsatz in EUR";`) {
			c.headerFound = true
		}
		return nil
	}

	handlers := []handler{handleLastschrif, handleWertpapiere, handleVisa, handleVisaMonthlyPayment}

	tokens := strings.Split(line, ";")
	length := len(tokens)
	if length < 5 {
		return nil
	}
	buchungsTag, _ := time.Parse(comdirectDateFormat, removeQuotes(tokens[0]))
	vorgang := removeQuotes(tokens[2])
	buchungsText := removeQuotes(tokens[3])
	umsatz := removeQuotes(tokens[4])
	if length == 6 {
		buchungsText = removeQuotes(tokens[4])
		umsatz = removeQuotes(tokens[5])
	}
	t := Transaction{}
	t.Date = buchungsTag
	t.Value = parseValue(umsatz)
	t.Comment = buchungsText
	t.Category = vorgang

	if t.Date.Before(c.fromDate) {
		return nil
	}

	for _, h := range handlers {
		if h(&t) {
			break
		}
	}

	c.sub.substitute(&t)
	filterRef(&t)

	return &t
}

func parseValue(str string) float64 {
	rawValue := strings.Replace(removeQuotes(str), ".", "", -1)
	value, _ := strconv.ParseFloat(strings.Replace(rawValue, ",", ".", 1), 64)
	return value
}

func handleLastschrif(t *Transaction) bool {
	if strings.Contains(t.Payee, "Lastschrift") || strings.Contains(t.Payee, "Überweisung") {
		s := strings.Split(t.Comment, " Buchungstext: ")
		t.Comment = s[1]
		t.Payee = strings.Replace(s[0], "Auftraggeber: ", "", 1)
		t.Payee = strings.Replace(t.Payee, "Empfänger: ", "", 1)

		return true
	}
	return false
}

func handleWertpapiere(t *Transaction) bool {
	if strings.Contains(t.Payee, "Wertpapiere") {
		t.Comment = strings.Replace(t.Comment, "Buchungstext: ", "", 1)
		t.Payee = "Transfer: .comdirect Depot"
		return true
	}
	return false
}

func handleVisa(t *Transaction) bool {
	if strings.Contains(t.Payee, "Visa-Umsatz") {
		t.Payee = t.Comment
		t.Comment = ""
		return true
	}
	return false
}

func handleVisaMonthlyPayment(t *Transaction) bool {
	if strings.Contains(t.Payee, "Visa-Kartenabrechnung") {
		t.Payee = "Transfer: .comdirect"
		t.Comment = ""
		return true
	}
	return false
}

func filterRef(t *Transaction) {
	t.Comment = strings.Split(t.Comment, " End-to-End-Ref.:")[0]
}

func removeQuotes(s string) string {
	return strings.Replace(s, "\"", "", -1)
}
