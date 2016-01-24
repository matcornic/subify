package opensubtitles

import (
	"errors"
	"path"

	"github.com/oz/osdb"
	logger "github.com/spf13/jwalterweatherman"
)

// API entry point
type API struct {
}

// Download downloads the OpenSubtitles subtitle from a video
func (s API) Download(videoPath string, language string) (subtitlePath string, err error) {
	logger.INFO.Println("Downloading subtitle with OpenSubtitles...")

	c, err := osdb.NewClient()
	if err != nil {
		return "", err
	}

	// Anonymous login
	if err = c.LogIn("", "", ""); err != nil {
		return "", err
	}
	languages := []string{language}

	// Search file
	subs, err := c.FileSearch(videoPath, languages)
	if err != nil {
		return "", err
	}

	// Keep best one
	best := subs.Best()
	if best == nil {
		return "", errors.New("Did not find best subtitle for this video")
	}

	// Saving to disk
	subtitlePath = videoPath[0:len(videoPath)-len(path.Ext(videoPath))] + ".srt"
	if err := c.DownloadTo(best, subtitlePath); err != nil {
		return "", err
	}

	logger.INFO.Println("Original name of subtitle :", best.SubFileName)

	return subtitlePath, nil
}

// Upload uploads the subtitle to OpenSubtitles, for the given video
func (s API) Upload(subtitlePath string, langauge string, videoPath string) error {
	return errors.New("Not yet implemented")
}
