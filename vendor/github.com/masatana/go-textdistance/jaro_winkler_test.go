package textdistance

import (
	"testing"
)

type jaroWinkerDistanceTest struct {
	s1       string
	s2       string
	expected float64
}

func (jdt *jaroWinkerDistanceTest) equals(dis float64) bool {
	if jdt.expected-1e-3 < dis && dis < jdt.expected+1e-3 {
		return true
	}
	return false
}

var JaroWinklerDistanceTests = []jaroWinkerDistanceTest{
	{
		"MARTHA",
		"MARHTA",
		0.961,
	},
	{
		"DWAYNE",
		"DUANE",
		0.84,
	},
	{
		"DIXON",
		"DICKSONX",
		0.813,
	},
	{
		"あいうえおか",
		"あいうおえか",
		0.961,
	},
}

func TestJaroWinklerDistance(t *testing.T) {
	for _, jwt := range JaroWinklerDistanceTests {
		actual := JaroWinklerDistance(jwt.s1, jwt.s2)
		if !jwt.equals(actual) {
			t.Errorf("Jaro-Winkler distance of %s and %s: want %f got %f", jwt.s1, jwt.s2, jwt.expected, actual)
		}
	}
}
