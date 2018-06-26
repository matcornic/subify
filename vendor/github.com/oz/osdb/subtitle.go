package osdb

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

// A Subtitle with its many OSDB attributes...
type Subtitle struct {
	IDMovie            string `xmlrpc:"IDMovie"`
	IDMovieImdb        string `xmlrpc:"IDMovieImdb"`
	IDSubMovieFile     string `xmlrpc:"IDSubMovieFile"`
	IDSubtitle         string `xmlrpc:"IDSubtitle"`
	IDSubtitleFile     string `xmlrpc:"IDSubtitleFile"`
	ISO639             string `xmlrpc:"ISO639"`
	LanguageName       string `xmlrpc:"LanguageName"`
	MatchedBy          string `xmlrpc:"MatchedBy"`
	MovieByteSize      string `xmlrpc:"MovieByteSize"`
	MovieFPS           string `xmlrpc:"MovieFPS"`
	MovieHash          string `xmlrpc:"MovieHash"`
	MovieImdbRating    string `xmlrpc:"MovieImdbRating"`
	MovieKind          string `xmlrpc:"MovieKind"`
	MovieName          string `xmlrpc:"MovieName"`
	MovieNameEng       string `xmlrpc:"MovieNameEng"`
	MovieReleaseName   string `xmlrpc:"MovieReleaseName"`
	MovieTimeMS        string `xmlrpc:"MovieTimeMS"`
	MovieYear          string `xmlrpc:"MovieYear"`
	MovieFileName      string `xmlrpc:"MovieName"`
	QueryNumber        string `xmlrpc:"QueryNumber"`
	SeriesEpisode      string `xmlrpc:"SeriesEpisode"`
	SeriesIMDBParent   string `xmlrpc:"SeriesIMDBParent"`
	SeriesSeason       string `xmlrpc:"SeriesSeason"`
	SubActualCD        string `xmlrpc:"SubActualCD"`
	SubAddDate         string `xmlrpc:"SubAddDate"`
	SubAuthorComment   string `xmlrpc:"SubAuthorComment"`
	SubBad             string `xmlrpc:"SubBad"`
	SubComments        string `xmlrpc:"SubComments"`
	SubDownloadLink    string `xmlrpc:"SubDownloadLink"`
	SubDownloadsCnt    string `xmlrpc:"SubDownloadsCnt"`
	SubFeatured        string `xmlrpc:"SubFeatured"`
	SubFileName        string `xmlrpc:"SubFileName"`
	SubFormat          string `xmlrpc:"SubFormat"`
	SubHash            string `xmlrpc:"SubHash"`
	SubHD              string `xmlrpc:"SubHD"`
	SubHearingImpaired string `xmlrpc:"SubHearingImpaired"`
	SubLanguageID      string `xmlrpc:"SubLanguageID"`
	SubRating          string `xmlrpc:"SubRating"`
	SubSize            string `xmlrpc:"SubSize"`
	SubSumCD           string `xmlrpc:"SubSumCD"`
	SubtitlesLink      string `xmlrpc:"SubtitlesLink"`
	UserID             string `xmlrpc:"UserID"`
	UserNickName       string `xmlrpc:"UserNickName"`
	UserRank           string `xmlrpc:"UserRank"`
	ZipDownloadLink    string `xmlrpc:"ZipDownloadLink"`
}

// A collection of subtitles
type Subtitles []Subtitle

// The best subtitle of the collection, for some definition of "best" at
// least.
func (subs Subtitles) Best() *Subtitle {
	if len(subs) > 0 {
		return &subs[0]
	}
	return nil
}

// SubtitleFile contains file data as returned by OSDB's API, that is to
// say: gzip-ped and base64-encoded text.
type SubtitleFile struct {
	Id     string `xmlrpc:"idsubtitlefile"`
	Data   string `xmlrpc:"data"`
	reader *gzip.Reader
}

// A Reader for the subtitle file contents (decoded, and decompressed).
func (sf *SubtitleFile) Reader() (r *gzip.Reader, err error) {
	if sf.reader != nil {
		return sf.reader, err
	}

	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(sf.Data))
	sf.reader, err = gzip.NewReader(dec)

	return sf.reader, err
}

// Build a Subtitle struct for a file, suitable for osdb.HasSubtitles()
func NewSubtitleWithFile(movie_file string, sub_file string) (s Subtitle, err error) {
	s.SubFileName = path.Base(sub_file)
	// Compute md5 sum
	sub_io, err := os.Open(sub_file)
	if err != nil {
		return
	}
	defer sub_io.Close()
	h := md5.New()
	_, err = io.Copy(h, sub_io)
	if err != nil {
		return
	}
	s.SubHash = fmt.Sprintf("%x", h.Sum(nil))

	// Movie filename, byte-size, & hash.
	s.MovieFileName = path.Base(movie_file)
	movie_io, err := os.Open(movie_file)
	if err != nil {
		return
	}
	defer movie_io.Close()
	stat, err := movie_io.Stat()
	if err != nil {
		return
	}
	s.MovieByteSize = strconv.FormatInt(stat.Size(), 10)
	movie_hash, err := HashFile(movie_io)
	if err != nil {
		return
	}
	s.MovieHash = fmt.Sprintf("%x", movie_hash)
	return
}

// Convert Subtitle to a map[string]string{}, because OSDB requires a
// specific structure to match subtitles when uploading (or trying to).
func (subs *Subtitles) toUploadParams() (map[string]interface{}, error) {
	subMap := map[string]interface{}{}
	for i, s := range *subs {
		key := "cd" + strconv.Itoa(i+1) // keys are cd1, cd2, ...
		param := map[string]string{
			"subhash":       s.SubHash,
			"subfilename":   s.SubFileName,
			"moviehash":     s.MovieHash,
			"moviebytesize": s.MovieByteSize,
			"moviefilename": s.MovieFileName,
		}
		subMap[key] = param
	}

	return subMap, nil
}
