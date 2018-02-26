package textdistance

import (
	"testing"
)

type jaccardSimilarityTest struct {
	s1       string
	s2       string
	expected float64
}

func (jst *jaccardSimilarityTest) equals(sim float64) bool {
	if jst.expected-1e-3 < sim && sim < jst.expected+1e-3 {
		return true
	}
	return false
}

var JaccardSimilarityTests = []jaccardSimilarityTest{
	{
		"c h",
		"a b c d e f g h",
		0.25,
	},
}

func TestJaccardSimilarity(t *testing.T) {
	for _, js := range JaccardSimilarityTests {
		actual := JaccardSimilarity(js.s1, js.s2, nil)
		if !js.equals(actual) {
			t.Errorf("Jaccard similarity of %s and %s: want %f got %f", js.s1, js.s2, js.expected, actual)
		}
	}
}
