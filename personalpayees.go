package main

type personalPayees struct {
}

func (personalPayees) substitute(payee string, memo string) (string, string) {
	for _, replacement := range replacements {
		if CaseInsensitiveContains(payee, replacement.in) {
			payee = replacement.out
			memo = ""
		}
	}
	return payee, memo
}

var replacements = []struct {
	in  string
	out string
}{
	{"E-CENTER", "Edeka"},
	{"EDEKA", "Edeka"},
	{"ROSSMANN", "Rossmann"},
	{"ALDI", "Aldi"},
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
	{"Amazon DE Marketplace", "Amazon Marketplace"},
}
