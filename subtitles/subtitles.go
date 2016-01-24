package subtitles

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/matcornic/subify/subtitles/subdb"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

// Subtitles defines the interface
type Subtitles interface {
	Download(videoPath string, language string) ([]byte, error)
	Upload(subtitlePath string, language string, videoPath string) error
}

// Download the subtitle from the video identified by its path
func Download(videoPath string) error {
	// APIs to download subtitles.
	apis := []Subtitles{subdb.API{}}
	var subtitle []byte
	var err error

	// Run through different APIs to get the subtitle. Stops when found
	for _, api := range apis {
		subtitle, err = api.Download(videoPath, viper.GetString("language"))
		if err == nil {
			break
		}
	}

	// Exit if no subtitle found, even after searching in all APIs
	if err != nil {
		return err
	}

	// Save the content to file
	subtitlePath := buildSubtitleName(videoPath)
	err = saveSubtitle(subtitle, subtitlePath)
	if err != nil {
		return err
	}

	logger.INFO.Println("Subtitle (", viper.GetString("language"), ") found and saved to ", subtitlePath)

	return nil
}

func buildSubtitleName(pathVideo string) string {
	var extension = filepath.Ext(pathVideo)
	var name = pathVideo[0 : len(pathVideo)-len(extension)]
	return name + ".srt"
}

func saveSubtitle(content []byte, subtitlePath string) error {
	err := ioutil.WriteFile(subtitlePath, content, 0644)
	if err != nil {
		return fmt.Errorf("Can't save the file %v because of : %v", subtitlePath, err)
	}
	return nil
}
