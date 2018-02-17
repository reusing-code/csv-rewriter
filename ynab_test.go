package csvrewrite

import "testing"

var formatValueTests = []struct {
	in  int
	out string
}{
	{100, "1.00"},
	{100000, "1000.00"},
	{100000000, "1000000.00"},
	{-567, "-5.67"},
	{-123456789, "-1234567.89"},
	{0, "0.00"},
	{-1, "-0.01"},
}

func TestFormatValue(t *testing.T) {
	for _, test := range formatValueTests {
		result := formatValue(test.in)
		if result != test.out {
			t.Errorf("Value '%d' formatted wrong. Expected %s was %s", test.in, test.out, result)
			continue
		}
	}
}
