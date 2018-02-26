package textdistance

import "testing"

type jaroDistanceTest struct {
	s1       string
	s2       string
	expected float64
}

func (jdt *jaroDistanceTest) equals(dis float64) bool {
	if jdt.expected-1e-3 < dis && dis < jdt.expected+1e-3 {
		return true
	}
	return false
}

var JaroDistanceTests = []jaroDistanceTest{
	{
		"MARTHA",
		"MARHTA",
		0.94444,
	},
	{
		"DWAYNE",
		"DUANE",
		0.822,
	},
	{
		"DIXON",
		"DICKSONX",
		0.767,
	},
	{
		"HELLO WORLD",
		"HELLO 我的世界",
		0.715,
	},
}

func TestJaroDistance(t *testing.T) {
	for _, jt := range JaroDistanceTests {
		actual, _ := JaroDistance(jt.s1, jt.s2)
		if !jt.equals(actual) {
			t.Errorf("Jaro distance of %s and %s: want %f got %f", jt.s1, jt.s2, jt.expected, actual)
		}
	}
}
