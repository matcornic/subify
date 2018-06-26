package textdistance

import (
	"github.com/deckarep/golang-set"
	"strings"
)

// JaccardSimilarity, as known as the Jaccard Index, compares the similarity of sample sets.
// This doesn't measure similarity between texts, but if regarding a text as bag-of-word,
// it can apply.
func JaccardSimilarity(s1, s2 string, f func(string) mapset.Set) float64 {
	if s1 == s2 {
		return 1.0
	}
	if f == nil {
		f = convertStringToSet
	}
	s1set := f(s1)
	s2set := f(s2)
	s1ands2 := s1set.Intersect(s2set).Cardinality()
	s1ors2 := s1set.Union(s2set).Cardinality()
	return float64(s1ands2) / float64(s1ors2)
}

func convertStringToSet(s string) mapset.Set {
	set := mapset.NewSet()
	for _, token := range strings.Fields(s) {
		set.Add(token)
	}
	return set
}
