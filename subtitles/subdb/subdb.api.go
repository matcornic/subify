package subdb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/matcornic/subify/subtitles/languages"
)

const (
	userAgent = "SubDB/1.0 (Subify/0.1; http://github.com/matcornic/subify)"
	devURL    = "http://sandbox.thesubd.com/"
	prodURL   = "http://api.thesubdb.com/"
)

// API entry point
type API struct {
	Name    string
	Aliases []string
}

// New creates a new API for OpenSubtitles
func New() API {
	return API{
		Name:    "SubDB",
		Aliases: []string{"subdb"},
	}
}

// Download downloads the SubDB subtitle from a video
func (s API) Download(videoPath string, language lang.Language) (subtitlePath string, err error) {
	// Get unique hash to identify video
	hash, err := getHashOfVideo(videoPath)
	if err != nil {
		return "", err
	}
	// Call SubDB API to get subtitle
	if language.SubDB == "" {
		return "", errors.New("Language exists but is not available for SubDB")
	}
	subtitle, err := subtitles(hash, language.SubDB)
	if err != nil {
		return "", err
	}

	// Save the content to file
	subtitlePath = videoPath[0:len(videoPath)-len(path.Ext(videoPath))] + ".srt"

	err = ioutil.WriteFile(subtitlePath, subtitle, 0644)
	if err != nil {
		return "", fmt.Errorf("Can't save the file %v because of : %v", subtitlePath, err)
	}

	return subtitlePath, nil
}

// Upload uploads the subtitle to SubDB, for the given video
func (s API) Upload(subtitlePath string, langauge lang.Language, videoPath string) error {
	return errors.New("Not yet implemented")
}

// GetName returns the name of the api
func (s API) GetName() string {
	return s.Name
}

// GetAliases returns aliases to identify this API
func (s API) GetAliases() []string {
	return s.Aliases
}
