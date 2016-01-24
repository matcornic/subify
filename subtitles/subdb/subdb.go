package subdb

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/lafikl/fluent"
	"github.com/matcornic/subify/common/config"
	"github.com/matcornic/subify/common/utils"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

const (
	userAgent       = "SubDB/1.0 (Subify/0.1; http://github.com/matcornic/subify)"
	devURL          = "http://sandbox.thesubdb.com/"
	prodURL         = "http://api.thesubdb.com/"
	defaultLanguage = "en"
)

// SubDBOptions describes parameters to the SubDB API
type options struct {
	Action   string `url:"action,omitempty"`
	Hash     string `url:"hash,omitempty"`
	Language string `url:"language,omitempty"`
}

func buildURL(hash string, language string) string {
	baseURL := prodURL
	if config.Dev {
		baseURL = devURL
	}
	opt := options{
		Action:   "download",
		Hash:     hash,
		Language: viper.GetString("language")}
	v, _ := query.Values(opt)

	url := baseURL + "?" + v.Encode()
	utils.VerbosePrintln(logger.INFO, "SubdbURL is : "+url)

	return url
}

// Subtitles get the subtitles from the hash of a video
func Subtitles(hash string) []byte {

	// Build request
	req := fluent.New()
	req.Get(buildURL(hash, defaultLanguage)).
		SetHeader("User-Agent", userAgent).
		InitialInterval(time.Duration(time.Millisecond)).
		Retry(3)

	// Execute the request
	res, err := req.Send()
	if err != nil {
		utils.ExitPrintError(err, "Can't reach the SubDB Web API. Are you connected to the Internet ?")
	}
	if res.StatusCode != 200 {
		utils.ExitVerbose(fmt.Sprintf("Response : %v", res),
			`Subtitle not found with SubDB Web API. Try with another language (See 'subify dl -h').
You may contribute to the community by updating their database (see 'subify upload -h')`)
	}

	// Extract the subtitles from the response
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		utils.ExitPrintError(err, "Can't read the content of the subtitles dowloaded from Subdb")
	}

	return content
}
