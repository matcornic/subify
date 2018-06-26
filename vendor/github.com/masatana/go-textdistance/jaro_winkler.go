package textdistance

// JaroWinklerDistance calculates jaro-winkler distance between s1 and s2.
// This implementation is influenced by an implementation of [lucene](http://lucene.apache.org/)
// Note that this calculation's result is normalized ( the result will be bewtwen 0 and 1)
// and if t1 and t2 are exactly the same, the result is 1.0.
func JaroWinklerDistance(s1, s2 string) float64 {
	threshold := 0.7
	jaroDistance, prefix := JaroDistance(s1, s2)
	if jaroDistance < threshold {
		return jaroDistance
	}
	return jaroDistance + (float64(Min(prefix, 4)) * 0.1 * (1 - jaroDistance))
}
