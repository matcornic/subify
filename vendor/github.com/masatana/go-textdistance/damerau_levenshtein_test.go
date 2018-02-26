package textdistance

import "testing"

type DamerauLevenshteinDistanceTest struct {
	s1       string
	s2       string
	expected int
}

var DamerauLevenshteinDistanceTests = []DamerauLevenshteinDistanceTest{
	{
		"ABC",
		"CDA",
		3,
	},
}

func TestDamerauLevenshteinDistance(t *testing.T) {
	for _, dlt := range DamerauLevenshteinDistanceTests {
		actual := DamerauLevenshteinDistance(dlt.s1, dlt.s2)
		if dlt.expected != actual {
			t.Errorf("DamerauLevenshtein distance of %s and %s: want %d got %d", dlt.s1, dlt.s2, dlt.expected, actual)
		}
	}
}
