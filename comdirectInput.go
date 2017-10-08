package main

import (
	"fmt"
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

func NewComdirectInput(fromDate time.Time, sub payeeSubstitution) *ComdirectInput {
	com := ComdirectInput{}
	com.headerFound = false
	com.fromDate = fromDate
	com.sub = sub
	return &com
}

func (c *ComdirectInput) processLine(line string) *Transaction {
	if !c.headerFound {
		if strings.EqualFold(line, `"Buchungstag";"Wertstellung (Valuta)";"Vorgang";"Buchungstext";"Umsatz in EUR";`) ||
			strings.EqualFold(line, `"Buchungstag";"Umsatztag";"Vorgang";"Referenz";"Buchungstext";"Umsatz in EUR";`) {
			c.headerFound = true
		}
		return nil
	}

	handlers := []handler{handleLastschrift, handleWertpapiere, handleVisa, handleVisaMonthlyPayment}

	tokens := splitLine(line, ';')
	length := len(tokens)
	if length < 5 {
		return nil
	}
	buchungsTag, _ := time.Parse(comdirectDateFormat, tokens[0])
	vorgang := tokens[2]
	buchungsText := tokens[3]
	umsatz := tokens[4]
	if length == 7 {
		buchungsText = tokens[4]
		umsatz = tokens[5]
	}
	var err error = nil
	t := Transaction{}
	t.Date = buchungsTag
	t.ValueCent, err = parseValue(umsatz)
	t.Comment = buchungsText
	t.Category = vorgang

	if err != nil {
		fmt.Errorf("error in parsing value from line '%s'", line)
		fmt.Println(err)
		return nil
	}

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

func parseValue(str string) (int, error) {
	rawValue := strings.Replace(str, ".", "", -1)
	rawValue = strings.Replace(rawValue, ",", "", -1)
	value, err := strconv.ParseInt(rawValue, 10, 32)
	return int(value), err
}

func handleLastschrift(t *Transaction) bool {
	if strings.Contains(t.Category, "Lastschrift") || strings.Contains(t.Category, "Überweisung") {
		s := strings.Split(t.Comment, " Buchungstext: ")
		t.Comment = s[1]
		t.Payee = strings.Replace(s[0], "Auftraggeber: ", "", 1)
		t.Payee = strings.Replace(t.Payee, "Empfänger: ", "", 1)
		t.Category = ""
		return true
	}
	return false
}

func handleWertpapiere(t *Transaction) bool {
	if strings.Contains(t.Category, "Wertpapiere") {
		t.Comment = strings.Replace(t.Comment, "Buchungstext: ", "", 1)
		t.Payee = "Transfer: .comdirect Depot"
		t.Category = ""
		return true
	}
	return false
}

func handleVisa(t *Transaction) bool {
	if strings.Contains(t.Category, "Visa-Umsatz") {
		t.Payee = t.Comment
		t.Comment = ""
		t.Category = ""
		return true
	}
	return false
}

func handleVisaMonthlyPayment(t *Transaction) bool {
	if strings.Contains(t.Category, "Visa-Kartenabrechnung") {
		t.Payee = "Transfer: .comdirect"
		t.Comment = ""
		t.Category = ""
		return true
	}
	return false
}

func filterRef(t *Transaction) {
	t.Comment = strings.Split(t.Comment, " End-to-End-Ref.:")[0]
}

func splitLine(s string, separator rune) []string {
	inQuotes := false
	var result = make([]string, 0)
	curStr := ""
	for _, runeValue := range s {
		if runeValue == '"' {
			inQuotes = !inQuotes
			continue
		}
		if inQuotes {
			curStr += string(runeValue)
			continue
		}
		if runeValue == separator {
			result = append(result, curStr)
			curStr = ""
			continue
		}
		curStr += string(runeValue)

	}
	result = append(result, curStr)
	return result
}
