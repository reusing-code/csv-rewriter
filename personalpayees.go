package csvrewrite

import "strings"

type PersonalPayees struct {
}

func (PersonalPayees) substitute(t *Transaction) {
	for _, replacement := range replacements {
		if CaseInsensitiveContains(t.Payee, replacement.in) {
			t.Payee = replacement.out
			t.Category = string(replacement.cat)
			t.Comment = ""
		}
	}
}

type Category string

const (
	GROCERIES      Category = "Everyday Expenses:Groceries"
	FUEL           Category = "Everyday Expenses:Fuel"
	SPENDING_MONEY Category = "Everyday Expenses:Spending Money"
	MEDICAL        Category = "Everyday Expenses:Medical"
	CLOTHING       Category = "Everyday Expenses:Clothing"
	DRUGSTORE      Category = "Everyday Expenses:Drogerieartikel"
	BARBER         Category = "Everyday Expenses:Barber"
	GAMES          Category = "Everyday Expenses:Games/DVDs/Software etc."
	CAR_MISC       Category = "Everyday Expenses:Car Miscellaneous"
	GIFTS          Category = "Everyday Expenses:Geschenke"
	TRAVELING      Category = "Everyday Expenses:Traveling"
	CINEMA         Category = "Everyday Expenses:Kino"
	RANDOM_CRAP    Category = "Everyday Expenses:Random Crap"
	TOOLS          Category = "Everyday Expenses:Werkzeug"
	TINKERING      Category = "Everyday Expenses:Gebastel"
	HOUSEHOLD      Category = "Everyday Expenses:Haushalt"
	BIKE           Category = "Everyday Expenses:Fahrrad"
	BOOKS          Category = "Everyday Expenses:Bücher"
	HARDWARE       Category = "Everyday Expenses:Hardware"
	HOBBIES        Category = "Everyday Expenses:Hobbies"
	RENT           Category = "Monthly Bills:Rent/Mortgage"
	PHONE          Category = "Monthly Bills:Phone"
	INTERNET       Category = "Monthly Bills:Internet"
	NETFLIX        Category = "Monthly Bills:Netflix etc."
	ELECTRICITY    Category = "Monthly Bills:Electricity"
	SPORT          Category = "Monthly Bills:Sport"
	CAR_LEASING    Category = "Monthly Bills:Auto Leasing"
	GIVING         Category = "Giving:Charitable"
	EMERGENCY      Category = "Rainy Day Funds:Emergency Fund"
	SAVING         Category = "Savings Goals:Saving (House, Car, Vacation)"
	GEZ            Category = "Sporadische Ausgaben:GEZ Verbrecher"
	HOSTING        Category = "Sporadische Ausgaben:Domain/Server/Etc"
	ETC            Category = "Sporadische Ausgaben:Sonstiges"
	INSURANCE      Category = "Sporadische Ausgaben:Versicherungen"
	PARTY          Category = "Sporadische Ausgaben:Feiern"
	TRAINING       Category = "Sporadische Ausgaben:Fortbildung"
	BANKING        Category = "Sporadische Ausgaben:Bank"
	NONE           Category = ""
)

var replacements = []struct {
	in  string
	out string
	cat Category
}{
	{"E-CENTER", "Edeka", GROCERIES},
	{"E CENTER", "Edeka", GROCERIES},
	{"EDEKA", "Edeka", GROCERIES},
	{"ALDI", "Aldi", GROCERIES},
	{"Lidl", "Lidl", GROCERIES},
	{"Norma", "Norma", GROCERIES},
	{"REWE", "REWE", GROCERIES},
	{"KAUFLAND", "Kaufland", GROCERIES},
	{"ROSSMANN", "Rossmann", DRUGSTORE},
	{"DM-Drogerie", "DM", DRUGSTORE},
	{"Bauhaus", "Bauhaus", HOUSEHOLD},
	{"Dehner", "Dehner", HOUSEHOLD},
	{"Willner", "Willner", BIKE},
	{"C+A", "C+A", CLOTHING},
	{"Schuh Muecke", "Schuh Mücke", CLOTHING},
	{"OUTLET", "Outlet", CLOTHING},
	{"TOM TAILOR", "Tom Tailor", CLOTHING},
	{"MCDONALDS", "McD", SPENDING_MONEY},
	{"AMAZON.DE", "Amazon", NONE},
	{"AMAZON EU S.A R.L.", "Amazon", NONE},
	{"Amazon DE Marketplace", "Amazon Marketplace", NONE},
	{"AMAZON PAYMENTS EUROPE S.C.A.", "Amazon Marketplace", NONE},
	{"E-TANKEN", "Tankstelle", FUEL},
	{"Walther", "Walther", FUEL},
	{"TOTAL", "TOTAL", FUEL},
	{"SHELL", "Shell", FUEL},
	{"ARAL", "Aral", FUEL},
	{"1u1", "1und1", HOSTING},
	{"Hetzner", "Hetzner", HOSTING},
	{"Netflix", "Netflix", NETFLIX},
	{"ERGO", "Ergo Direkt", INSURANCE},
	{"COSMOS", "Transfer: Cosmos Leben", SAVING},
	{"Anna Roiser", "Anna Roiser", RENT},
	{"COM-IN", "COM-IN", INTERNET},
	{"Stadtwerke Ingolstadt", "Stadtwerke Ingolstadt", ELECTRICITY},
	{"Visa-Monatsabrechnung", "Transfer: Visa", NONE},
	{"BARCLAYCARD", "Transfer: Barclaycard", NONE},
}

func CaseInsensitiveContains(s, substr string) bool {
	s, substr = strings.ToUpper(s), strings.ToUpper(substr)
	return strings.Contains(s, substr)
}
