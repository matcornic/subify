package textdistance

import "strings"

// LevenshteinDistance calculates the levenshtein distance between s1 and s2.
// Reference: [Levenshtein Distance](http://en.wikipedia.org/wiki/Levenshtein_distance)
// Note that this calculation's result isn't normalized. (not between 0 and 1.)
// and if s1 and s2 are exactly the same, the result is 0.
func LevenshteinDistance(s1, s2 string) int {
	if s1 == s2 {
		return 0
	}
	s1Array := strings.Split(s1, "")
	s2Array := strings.Split(s2, "")
	lenS1Array := len(s1Array)
	lenS2Array := len(s2Array)
	if lenS1Array == 0 {
		return lenS2Array
	}
	if lenS2Array == 0 {
		return lenS1Array
	}
	m := make([][]int, lenS1Array+1)
	for i := range m {
		m[i] = make([]int, lenS2Array+1)
	}
	for i := 0; i < lenS1Array+1; i++ {
		for j := 0; j < lenS2Array+1; j++ {
			if i == 0 {
				m[i][j] = j
			} else if j == 0 {
				m[i][j] = i
			} else {
				if s1Array[i-1] == s2Array[j-1] {
					m[i][j] = m[i-1][j-1]
				} else {
					m[i][j] = Min(m[i-1][j]+1, m[i][j-1]+1, m[i-1][j-1]+1)
				}
			}
		}
	}
	return m[lenS1Array][lenS2Array]
}
