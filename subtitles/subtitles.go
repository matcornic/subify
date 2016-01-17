package subtitles

import (
	"fmt"
	"github.com/matcornic/subify/common/config"
	"github.com/matcornic/subify/common/utils"
	"github.com/matcornic/subify/subtitles/subdb"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Download the subtitle from the video identified by its path
func Download(videoPath string) {
	// Get unique hash to identify video
	hash := utils.GetHashOfVideo(videoPath)
	// Call SubDB API to get subtitle
	subtitle := subdb.Subtitles(hash)
	// Save the content to file
	subtitlePath := buildSubtitleName(videoPath)
	saveSubtitle(subtitle, subtitlePath)

	fmt.Println("Subtitle (" + viper.GetString("language") + ") found and saved to " + subtitlePath)
}

func buildSubtitleName(pathVideo string) string {
	var extension = filepath.Ext(pathVideo)
	var name = pathVideo[0 : len(pathVideo)-len(extension)]
	return name + ".srt"
}

func saveSubtitle(content []byte, subtitlePath string) {
	err := ioutil.WriteFile(subtitlePath, content, 0644)
	if err != nil {
		if config.Verbose {
			fmt.Fprintln(os.Stderr, err)
		}
		utils.Exit("Can't save the file %v", subtitlePath)
	}
}
