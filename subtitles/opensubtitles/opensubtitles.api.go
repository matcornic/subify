package opensubtitles

import (
	"errors"
	"path"

	"github.com/matcornic/subify/subtitles/languages"
	"github.com/oz/osdb"
	logger "github.com/spf13/jwalterweatherman"
)

const (
	userAgent = "Subify 0.1.0"
)

// API entry point
type API struct {
	Name    string
	Aliases []string
}

// New creates a new API for OpenSubtitles
func New() API {
	return API{
		Name:    "OpenSubtitles",
		Aliases: []string{"os", "opensubtitles", "opensubtitle", "opensub", "osub", "osdb", "open"},
	}
}

// Download downloads the OpenSubtitles subtitle from a video
func (s API) Download(videoPath string, language lang.Language) (subtitlePath string, err error) {
	c, err := osdb.NewClient()
	if err != nil {
		return "", err
	}
	c.UserAgent = userAgent

	// Anonymous login
	if err = c.LogIn("", "", ""); err != nil {
		return "", err
	}
	if language.OpenSubtitles == "" {
		return "", errors.New("Language exists but is not available for OpenSubtitles")
	}
	languages := []string{language.OpenSubtitles}

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
