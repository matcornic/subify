package subtitles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitAPIsShouldExists(t *testing.T) {
	apis := InitAPIs([]string{"subdb", "oS"})
	assert.Equal(t, len(apis), 2, "Should have two apis")
	assert.Equal(t, apis[0].GetName(), "SubDB", "Should be SubDB")
	assert.Equal(t, apis[1].GetName(), "OpenSubtitles", "Should be OpenSubtitles")
}

func TestInitAPIsShouldExistsWithOneThatDoesNotExist(t *testing.T) {
	apis := InitAPIs([]string{"subDB", "os", "dontexist"})
	assert.Equal(t, len(apis), 2, "Should have two apis")
	assert.Equal(t, apis[0].GetName(), "SubDB", "Should be SubDB")
	assert.Equal(t, apis[1].GetName(), "OpenSubtitles", "Should be OpenSubtitles")
}
