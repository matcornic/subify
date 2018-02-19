package subtitles

import (
	"errors"
	"path"

	"github.com/oz/osdb"
	logger "github.com/spf13/jwalterweatherman"
)

const (
	osdbUserAgent = "Subify 0.1.0"
)

var osLangs = map[string]string{
	"afr": "afr",
	"alb": "alb",
	"ara": "ara",
	"arm": "arm",
	"bel": "bel",
	"ben": "ben",
	"bos": "bos",
	"bul": "bul",
	"bur": "bur",
	"cat": "cat",
	"chi": "chi",
	"cze": "cze",
	"dan": "dan",
	"dut": "dut",
	"eng": "eng",
	"epo": "epo",
	"est": "est",
	"fin": "fin",
	"fre": "fre",
	"geo": "geo",
	"ger": "ger",
	"glg": "glg",
	"ell": "ell",
	"heb": "heb",
	"hin": "hin",
	"hrv": "hrv",
	"hun": "hun",
	"ice": "ice",
	"ind": "ind",
	"ita": "ita",
	"jpn": "jpn",
	"kaz": "kaz",
	"kor": "kor",
	"lav": "lav",
	"ltz": "ltz",
	"lit": "lit",
	"mac": "mac",
	"mni": "mni",
	"mon": "mon",
	"nor": "nor",
	"per": "per",
	"pol": "pol",
	"por": "por",
	"rus": "rus",
	"scc": "scc",
	"sin": "sin",
	"slo": "slo",
	"slv": "slv",
	"spa": "spa",
	"swa": "swa",
	"swe": "swe",
	"syr": "syr",
	"tam": "tam",
	"tha": "tha",
	"tur": "tur",
	"ukr": "ukr",
	"vie": "vie",
	"rum": "rum",
	"pob": "pob",
	"zht": "zht",
	"zhe": "zhe",
}

// OSDBAPI entry point
type OSDBAPI struct {
	Name    string
	Aliases []string
}

// OpenSubtitles creates a new API for OpenSubtitles
func OpenSubtitles() OSDBAPI {
	return OSDBAPI{
		Name:    "OpenSubtitles",
		Aliases: []string{"os", "opensubtitles", "opensubtitle", "opensub", "osub", "osdb", "open"},
	}
}

// Download downloads the OpenSubtitles subtitle from a video
func (s OSDBAPI) Download(videoPath string, language Language) (subtitlePath string, err error) {
	c, err := osdb.NewClient()
	if err != nil {
		return "", err
	}
	c.UserAgent = osdbUserAgent

	// Anonymous login
	if err = c.LogIn("", "", ""); err != nil {
		return "", err
	}
	lang, ok := osLangs[language.ID]
	if !ok {
		return "", errors.New("Language exists but is not available for OpenSubtitles")
	}
	languages := []string{lang}

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
func (s OSDBAPI) Upload(subtitlePath string, langauge Language, videoPath string) error {
	return errors.New("Not yet implemented")
}

// GetName returns the name of the api
func (s OSDBAPI) GetName() string {
	return s.Name
}

// GetAliases returns aliases to identify this API
func (s OSDBAPI) GetAliases() []string {
	return s.Aliases
}
