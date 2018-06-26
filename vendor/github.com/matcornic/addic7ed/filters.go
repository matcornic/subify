package addic7ed

import (
	"regexp"
	"strings"
)

// WithLanguage is a filter first-class function, used to keep subtitle with given language
func WithLanguage(lang string) func(s Subtitle) bool {
	return func(s Subtitle) bool {
		return strings.EqualFold(strings.TrimSpace(s.Language), strings.TrimSpace(lang))
	}
}

// WithVersion is a filter first-class function, used to keep subtitle with given subtitle version
func WithVersion(version string) func(s Subtitle) bool {
	return func(s Subtitle) bool {
		return strings.EqualFold(strings.TrimSpace(s.Version), strings.TrimSpace(version))
	}
}

// WithVersionRegexp is a filter first-class function, used to keep subtitle with given subtitle version identified by a regex
func WithVersionRegexp(version *regexp.Regexp) func(s Subtitle) bool {
	return func(s Subtitle) bool {
		return version.MatchString(strings.TrimSpace(s.Version))
	}
}
