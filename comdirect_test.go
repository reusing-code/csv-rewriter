package csvrewrite

import "testing"

var splitTests = []struct {
	in  string
	out []string
}{
	{"a,b,c,d", []string{"a", "b", "c", "d"}},
	{`"a","b","c","d"`, []string{"a", "b", "c", "d"}},
	{`"a,b","b","c,abc,d","d,,,e"`, []string{"a,b", "b", "c,abc,d", "d,,,e"}},
	{"", []string{""}},
	{",", []string{"", ""}},
	{",,,", []string{"", "", "", ""}},
	{`'a','b','c','d'`, []string{"'a'", "'b'", "'c'", "'d'"}},
}

func TestLineSplitter(t *testing.T) {
	for _, test := range splitTests {
		result := splitLine(test.in, ',')
		if len(result) != len(test.out) {
			t.Errorf("Line '%s' splitted wrong. Len expected %d was %d", test.in, len(test.out), len(result))
			continue
		}

		for i := range test.out {
			if test.out[i] != result[i] {
				t.Errorf("Line '%s' splitted wrong. Token expected '%s' was '%s'", test.in, test.out[i], result[i])
			}
		}
	}
}

var parseValueTests = []struct {
	in  string
	out int
}{
	{"1,00", 100},
	{"1.00", 100},
	{"1,000.00", 100000},
	{"1.000.000,00", 100000000},
	{"-5,67", -567},
	{"-1,234,567.89", -123456789},
	{"0,00", 0},
	{"-0.01", -1},
}

func TestParseValue(t *testing.T) {
	for _, test := range parseValueTests {
		result, _ := parseValue(test.in)
		if result != test.out {
			t.Errorf("Value '%s' parsed wrong. Expected %d was %d", test.in, test.out, result)
			continue
		}
	}
}

var parseValueTestsInvalid = []struct {
	in string
}{
	{"a1,00"},
	{"1.00a"},
	{"1a1"},
	{"1;"},
	{":1"},
}

func TestParseValueInvalid(t *testing.T) {
	for _, test := range parseValueTestsInvalid {
		_, err := parseValue(test.in)
		if err == nil {
			t.Errorf("Invalid value '%s' parsed wrong. Expected error", test.in)
			continue
		}
	}
}
