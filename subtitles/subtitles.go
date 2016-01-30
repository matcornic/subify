package subtitles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/matcornic/subify/subtitles/languages"
	"github.com/matcornic/subify/subtitles/opensubtitles"
	"github.com/matcornic/subify/subtitles/subdb"
	logger "github.com/spf13/jwalterweatherman"
)

// Client defines the interface to get subtitles from API
type Client interface {
	Download(videoPath string, language lang.Language) (subtitlePath string, err error)
	Upload(subtitlePath string, language lang.Language, videoPath string) error
	GetName() string
}

// Download the subtitle from the video identified by its path
func Download(videoPath string, languages []string) error {
	// APIs to download subtitles.
	apis := []Client{
		subdb.API{},
		opensubtitles.API{},
	}
	var subtitlePath string
	var err error

	// Check languages
	l := lang.Languages.GetLanguages(languages)
	if len(l) == 0 {
		logger.ERROR.Println("Languages", languages, "are not available. Pick one ore more from the table below :")
		lang.Languages.Print(false)
		return fmt.Errorf("No languages is available for given languages : %v", languages)
	} else if len(languages) != len(l) {
		logger.WARN.Println("Some languages are not recognized. Given:", languages, "Found:", l.GetDescriptions())
	}

	// Run through languages
browselang:
	for i, lang := range l {
		// Run through different APIs to get the subtitle. Stops when found
		logger.INFO.Println("===> ("+strconv.Itoa(i+1)+") Searching subtitles for", lang.Description, "language")
		for j, api := range apis {
			logger.INFO.Println("=> (" + strconv.Itoa(i+1) + "." + strconv.Itoa(j+1) + ") Downloading subtitle with " + api.GetName() + "...")
			subtitlePath, err = api.Download(videoPath, lang)
			if err == nil {
				logger.INFO.Println(lang.Description, "subtitle found and saved to ", subtitlePath)
				break browselang
			} else {
				logger.INFO.Println("Subtitle not found because :", err.Error())
			}
			if (j + 1) < len(apis) {
				logger.INFO.Println("Trying with another API...")
			}
		}
		if err != nil {
			logger.INFO.Println("=> No subtitle found in", lang.Description, "language.")
		}
		if (i + 1) < len(l) {
			logger.INFO.Println("Trying with another language...")
		}
	}

	if err != nil {
		return fmt.Errorf("No %v subtitle found, even after searching in all APIs for all given languages", strings.Join(l.GetDescriptions(), ", nor "))
	}

	return nil
}
