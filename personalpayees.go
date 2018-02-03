package main

import "strings"

type personalPayees struct {
}

func (personalPayees) substitute(t *Transaction) {
	for _, replacement := range replacements {
		if CaseInsensitiveContains(t.Payee, replacement.in) {
			t.Payee = replacement.out
			t.Comment = ""
		}
	}
}

var replacements = []struct {
	in  string
	out string
}{
	{"E-CENTER", "Edeka"},
	{"EDEKA", "Edeka"},
	{"ROSSMANN", "Rossmann"},
	{"ALDI", "Aldi"},
	{"Lidl", "Lidl"},
	{"ERGO", "Ergo Direkt"},
	{"COSMOS", "Transfer: Cosmos Leben"},
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
	{"Visa-Monatsabrechnung", "Transfer: Visa"},
	{"BARCLAYCARD", "Barcleycard"},
	{"AMAZON.DE", "Amazon"},
	{"AMAZON EU S.A R.L.", "Amazon"},
	{"Amazon DE Marketplace", "Amazon Marketplace"},
	{"AMAZON PAYMENTS EUROPE S.C.A.", "Amazon Marketplace"},
	{"Walther", "Walther"},
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
