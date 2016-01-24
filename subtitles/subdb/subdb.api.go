package subdb

import "errors"

const (
	userAgent = "SubDB/1.0 (Subify/0.1; http://github.com/matcornic/subify)"
	devURL    = "http://sandbox.thesubd.com/"
	prodURL   = "http://api.thesubdb.com/"
)

// API entry point
type API struct {
}

// Download downloads the SubDB subtitle from a video
func (s API) Download(videoPath string, language string) ([]byte, error) {
	// Get unique hash to identify video
	hash, err := getHashOfVideo(videoPath)
	if err != nil {
		return []byte{}, err
	}
	// Call SubDB API to get subtitle
	subtitle, err := subtitles(hash, language)
	if err != nil {
		return []byte{}, err
	}

	return subtitle, nil
}

// Upload uploads the subtitle to SubDB, for the given video
func (s API) Upload(subtitlePath string, langauge string, videoPath string) error {
	return errors.New("Not yet implemented")
}
