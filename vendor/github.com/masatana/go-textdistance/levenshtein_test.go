package textdistance

import "testing"

type LevenshteinDistanceTest struct {
	s1       string
	s2       string
	expected int
}

var LevenshteinDistanceTests = []LevenshteinDistanceTest{
	{
		"kitten",
		"sitting",
		3,
	},
	{
		"book",
		"back",
		2,
	},
	{
		"test",
		"",
		4,
	},
	{
		"",
		"test",
		4,
	},
	{
		"日本語も",
		"",
		4,
	},
	{
		"",
		"大丈夫",
		3,
	},
	{
		"あaaあ",
		"あeoあ",
		2,
	},
}

func TestLevenshteinDistance(t *testing.T) {
	for _, lt := range LevenshteinDistanceTests {
		actual := LevenshteinDistance(lt.s1, lt.s2)
		if lt.expected != actual {
			t.Errorf("Levenshtein distance of %s and %s: want %d got %d", lt.s1, lt.s2, lt.expected, actual)
		}
	}
}
