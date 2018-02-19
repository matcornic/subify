package subtitles

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/lafikl/fluent"
	"github.com/matcornic/subify/common/config"
)

const (
	subDbUserAgent = "SubDB/1.0 (Subify/0.1; http://github.com/matcornic/subify)"
	subdbDevURL    = "http://sandbox.thesubdb.com/"
	subdbProdURL   = "http://api.thesubdb.com/"
)

var subdbLangs = map[string]string{
	"dut": "nl",
	"eng": "en",
	"fre": "fr",
	"ita": "it",
	"pol": "pl",
	"spa": "es",
	"swe": "sv",
	"tur": "tr",
	"rum": "ro",
	"pob": "pt",
}

// API entry point
type SubDBAPI struct {
	Name    string
	Aliases []string
}

// New creates a new API for OpenSubtitles
func SubDB() SubDBAPI {
	return SubDBAPI{
		Name:    "SubDB",
		Aliases: []string{"subdb"},
	}
}

// Download downloads the SubDB subtitle from a video
func (s SubDBAPI) Download(videoPath string, language Language) (subtitlePath string, err error) {
	// Get unique hash to identify video
	hash, err := getHashOfVideo(videoPath)
	if err != nil {
		return "", err
	}
	// Call SubDB API to get subtitle
	lang, ok := subdbLangs[language.ID]
	if !ok {
		return "", errors.New("Language exists but is not available for SubDB")
	}
	subtitle, err := subtitles(hash, lang)
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
func (s SubDBAPI) Upload(subtitlePath string, langauge Language, videoPath string) error {
	return errors.New("Not yet implemented")
}

// GetName returns the name of the api
func (s SubDBAPI) GetName() string {
	return s.Name
}

// GetAliases returns aliases to identify this API
func (s SubDBAPI) GetAliases() []string {
	return s.Aliases
}

// options describes parameters to the SubDB API
type options struct {
	Action   string `url:"action,omitempty"`
	Hash     string `url:"hash,omitempty"`
	Language string `url:"language,omitempty"`
}

//getHashOfVideo gets the hash used by SubDb to identify a video. Absolutely needed either to download or upload subtitles.
//The hash is composed by taking the first and the last 64kb of the video file, putting all together and generating a md5 of the resulting data (128kb).
func getHashOfVideo(filename string) (string, error) {
	readsize := 64 * 1024 // 64kb

	// Open Video
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Can't open file %v because of : %v ", filename, err.Error())
	}
	defer file.Close()

	// Get stats of file
	fi, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("Can't get stats for file %v because of : %v", filename, err.Error())
	}

	// Fill a buffer with first bytes of file
	bufB := make([]byte, readsize)
	_, err = file.Read(bufB)
	if err != nil {
		return "", fmt.Errorf("Can't read content of file %v because of : %v", filename, err.Error())
	}

	//Fill a buffer with last bytes of file
	bufE := make([]byte, readsize)
	n, err := file.ReadAt(bufE, fi.Size()-int64(len(bufE)))
	if err != nil {
		return "", fmt.Errorf("File is probably too small, can't read content of file %v because of : %v", filename, err.Error())
	}
	bufE = bufE[:n]

	// Generates MD5 of both bytes chain
	bufB = append(bufB, bufE...)
	hash := fmt.Sprintf("%x", md5.Sum(bufB))

	return hash, nil
}

func buildURL(hash string, language string) string {
	baseURL := subdbProdURL
	if config.Dev {
		fmt.Println("Dev mode")
		baseURL = subdbDevURL
	} else {
		fmt.Println("Prod mode")
	}
	opt := options{
		Action:   "download",
		Hash:     hash,
		Language: language}
	v, _ := query.Values(opt)

	url := baseURL + "?" + v.Encode()

	return url
}

// Subtitles get the subtitles from the hash of a video
func subtitles(hash string, language string) ([]byte, error) {

	// Build request
	req := fluent.New()
	req.Get(buildURL(hash, language)).
		SetHeader("User-Agent", subDbUserAgent).
		InitialInterval(time.Duration(time.Millisecond)).
		Retry(3)

	// Execute the request
	res, err := req.Send()
	if err != nil {
		return []byte{}, fmt.Errorf("Can't reach the SubDB Web API. Are you connected to the Internet ? %v", err.Error())
	}
	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf(`Subtitle not stored by SubDB`)
	}

	// Extract the subtitles from the response
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("The content of the subtitles dowloaded from Subdb is corrupted")
	}

	return content, nil
}
