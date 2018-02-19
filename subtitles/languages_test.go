package subtitles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLanguageShouldExists(t *testing.T) {
	language := Languages.GetLanguage("en")
	assert.NotNil(t, language)
	assert.Equal(t, language.Description, "English", "Should be english")
}

func TestGetLanguageShouldExistsWithDifferentCase(t *testing.T) {
	language := Languages.GetLanguage("En")
	assert.NotNil(t, language)
	assert.Equal(t, language.Description, "English", "Should be english")
}

func TestGetLanguageShouldExistsWithAlias(t *testing.T) {
	language := Languages.GetLanguage("eNg")
	assert.NotNil(t, language)
	assert.Equal(t, language.Description, "English", "Should be english")
}

func TestGetLanguageShouldBeNilWhenEmpty(t *testing.T) {
	language := Languages.GetLanguage("")
	assert.Nil(t, language)
}

func TestGetLanguageShouldBeNilWhenNotExist(t *testing.T) {
	language := Languages.GetLanguage("elfic")
	assert.Nil(t, language)
}

func TestGetLanguagesShouldExists(t *testing.T) {
	languages := Languages.GetLanguages([]string{"en", "fr"})
	assert.Equal(t, len(languages), 2, "Should have two languages")
	assert.Equal(t, languages[0].Description, "English", "Should be english")
	assert.Equal(t, languages[1].Description, "French", "Should be french")
}

func TestGetLanguagesShouldPartiallyExists(t *testing.T) {
	languages := Languages.GetLanguages([]string{"EN", "elfic"})
	assert.Equal(t, len(languages), 1, "Should have one language")
	assert.Equal(t, languages[0].Description, "English", "Should be english")
}
