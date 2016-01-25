package subtitles

import (
	"fmt"
	"strings"

	"github.com/matcornic/subify/subtitles/languages"
	"github.com/matcornic/subify/subtitles/opensubtitles"
	"github.com/matcornic/subify/subtitles/subdb"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Client defines the interface to get subtitles from API
type Client interface {
	Download(videoPath string, language lang.Language) (subtitlePath string, err error)
	Upload(subtitlePath string, language lang.Language, videoPath string) error
}

// Download the subtitle from the video identified by its path
func Download(videoPath string) error {
	// APIs to download subtitles.
	apis := []Client{
		subdb.API{},
		opensubtitles.API{},
	}
	var subtitlePath string
	var err error

	// Check language
	l := lang.Languages.GetLanguage(strings.ToLower(viper.GetString("language")))
	if l == nil {
		logger.ERROR.Println("Language", viper.GetString("language"), "is not available. Pick one from the table below :")
		lang.Languages.Print(false)
		return fmt.Errorf("Language %v is not available", viper.GetString("language"))
	}

	// Run through different APIs to get the subtitle. Stops when found
	for _, api := range apis {
		subtitlePath, err = api.Download(videoPath, *l)
		if err == nil {
			break
		} else {
			logger.INFO.Println("Subtitle not found because :", err.Error(), ". Trying with another API...")
		}
	}

	if err != nil {
		return fmt.Errorf("No %v subtitle found at all, even after searching in all APIs", l.Description)
	}

	logger.INFO.Println(l.Description, "subtitle found and saved to ", subtitlePath)

	return nil
}
