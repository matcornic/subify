package subtitles

import (
	"github.com/matcornic/subify/subtitles/opensubtitles"
	"github.com/matcornic/subify/subtitles/subdb"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Client defines the interface to get subtitles from API
type Client interface {
	Download(videoPath string, language string) (subtitlePath string, err error)
	Upload(subtitlePath string, language string, videoPath string) error
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

	// Run through different APIs to get the subtitle. Stops when found
	for _, api := range apis {
		subtitlePath, err = api.Download(videoPath, viper.GetString("language"))
		if err == nil {
			break
		} else {
			logger.INFO.Println("Subtitle not found. Trying with another API...")
		}
	}

	if err != nil {
		logger.ERROR.Println("No subtitle found at all, even after searching in all APIs")
		return err
	}

	logger.INFO.Println("Subtitle (", viper.GetString("language"), ") found and saved to ", subtitlePath)

	return nil
}
