package subdb

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/lafikl/fluent"
	"github.com/matcornic/subify/common/config"
	"github.com/matcornic/subify/common/utils"
	logger "github.com/spf13/jwalterweatherman"
)

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

	utils.VerbosePrintln(logger.INFO, "SubDB Hash of video is "+hash)

	return hash, nil
}

func buildURL(hash string, language string) string {
	baseURL := prodURL
	if config.Dev {
		baseURL = devURL
	}
	opt := options{
		Action:   "download",
		Hash:     hash,
		Language: language}
	v, _ := query.Values(opt)

	url := baseURL + "?" + v.Encode()
	utils.VerbosePrintln(logger.INFO, "SubdbURL is : "+url)

	return url
}

// Subtitles get the subtitles from the hash of a video
func subtitles(hash string, language string) ([]byte, error) {

	// Build request
	req := fluent.New()
	req.Get(buildURL(hash, language)).
		SetHeader("User-Agent", userAgent).
		InitialInterval(time.Duration(time.Millisecond)).
		Retry(3)

	// Execute the request
	res, err := req.Send()
	if err != nil {
		return []byte{}, fmt.Errorf("Can't reach the SubDB Web API. Are you connected to the Internet ? %v", err.Error())
	}
	if res.StatusCode != 200 {
		return []byte{}, fmt.Errorf(fmt.Sprintf("Response : %v", res),
			`Subtitle not found with SubDB Web API. Try with another language (See 'subify dl -h').
You may contribute to the community by updating their database (see 'subify upload -h')`)
	}

	// Extract the subtitles from the response
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Can't read the content of the subtitles dowloaded from Subdb")
	}

	return content, nil
}
