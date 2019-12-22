package subtitles

import (
	"errors"
	"path"

	"github.com/matcornic/addic7ed"
	logger "github.com/spf13/jwalterweatherman"
)

// Addic7edAPI is the endpoint for downloading Addic7ed subtitles
type Addic7edAPI struct {
	Name    string
	Aliases []string
}

var addic7edLangs = map[string]string{
	"alb": "Albanian",
	"ara": "Arabic",
	"arm": "Armenian",
	"aze": "Azerbaijani",
	"ben": "Bengali",
	"bos": "Bosnian",
	"bul": "Bulgarian",
	"cat": "Català",
	"chi": "Chinese (Simplified)",
	"hrv": "Croatian",
	"cze": "Czech",
	"dan": "Danish",
	"dut": "Dutch",
	"eng": "English",
	"est": "Estonian",
	"baq": "Euskera",
	"fin": "Finnish",
	"fre": "French",
	"frc": "French (Canadian)",
	"glg": "Galego",
	"ger": "German",
	"ell": "Greek",
	"heb": "Hebrew",
	"hin": "Hindi",
	"hun": "Hungarian",
	"ice": "Icelandic",
	"ind": "Indonesian",
	"ita": "Italian",
	"jpn": "Japanese",
	"tlh": "Klingon",
	"kor": "Korean",
	"lav": "Latvian",
	"lit": "Lithuanian",
	"mac": "Macedonian",
	"mal": "Malay",
	"nor": "Norwegian",
	"per": "Persian",
	"pol": "Polish",
	"por": "Portuguese",
	"pob": "Portuguese (Brazilian)",
	"rum": "Romanian",
	"rus": "Russian",
	"scc": "Serbian (Cyrillic)",
	"sin": "Sinhala",
	"slo": "Slovak",
	"slv": "Slovenian",
	"es":  "Spanish",
	"swe": "Swedish",
	"tam": "Tamil",
	"tha": "Thai",
	"tur": "Turkish",
	"ukr": "Ukrainian",
	"vie": "Vietnamese",
	"zht": "Chinese (Traditional)",
}

// Addic7ed creates a new API for Addic7ed
func Addic7ed() Addic7edAPI {
	return Addic7edAPI{
		Name:    "Addic7ed",
		Aliases: []string{"addic7ed", "ad7", "ad", "addicted", "addict", "add"},
	}
}

// Download downloads the Addic7ed subtitle from a video
func (s Addic7edAPI) Download(videoPath string, language Language) (subtitlePath string, err error) {
	c := addic7ed.New()

	lang, ok := addic7edLangs[language.ID]
	if !ok {
		return "", errors.New("Language exists but is not available for Addic7ed")
	}

	show, subtitle, err := c.SearchBest(path.Base(videoPath), lang)
	if err != nil {
		return "", err
	}

	// Saving to disk
	subtitlePath = videoPath[0:len(videoPath)-len(path.Ext(videoPath))] + "." + lang + ".srt"
	if err := subtitle.DownloadTo(subtitlePath); err != nil {
		return "", err
	}

	logger.INFO.Println("Original name of subtitle :", show)

	return subtitlePath, nil
}

// Upload uploads the subtitle to OpenSubtitles, for the given video
func (s Addic7edAPI) Upload(subtitlePath string, language Language, videoPath string) error {
	return errors.New("Not yet implemented")
}

// GetName returns the name of the api
func (s Addic7edAPI) GetName() string {
	return s.Name
}

// GetAliases returns aliases to identify this API
func (s Addic7edAPI) GetAliases() []string {
	return s.Aliases
}
