package addic7ed

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	textdistance "github.com/masatana/go-textdistance"
)

const userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:12.0) Gecko/20100101 Firefox/12.0"

// Client is the addic7ed client
type Client struct {
	// doc is the indexed document, representing the page
	doc *goquery.Document
}

// New creates an Addic7ed client, ready to interact with.
func New() *Client {
	return &Client{}
}

func (c *Client) findShowName() (string, error) {
	var show string
	c.doc.Find(".titulo").Contents().EachWithBreak(func(i int, s *goquery.Selection) bool {
		if !s.Is("small") {
			show = strings.TrimSpace(s.Text())
			return false
		}
		return true
	})
	if show == "" {
		return "", errors.New("not found")
	}
	return show, nil
}

func (c *Client) findResults() []string {
	results := []string{}
	c.doc.Find(".tabel").Each(func(i int, s *goquery.Selection) {
		s.Find("a").Each(func(j int, ss *goquery.Selection) {
			if url, ok := ss.Attr("href"); ok {
				results = append(results, url)
			}
		})
	})
	return results
}

func createDocFromURL(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// Avoid getting cached pages
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to reach addic7ed server: %v", err)
	}

	// We use goquery to fetch the page from Addic7ed in way that we can find data quickly like the JQuery way
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("Unable to construct document from server response: %v", err)
	}

	return doc, nil
}

// fetchShowPage get the addic7ed show page from Addic7ed website
// It uses search function of the website to get the page
// Return an error if the page is not found
// If more than one result is returned, we get the first one to match
func (c *Client) fetchShowPage(fileName string) (string, error) {

	doc, err := createDocFromURL(fmt.Sprintf("http://www.addic7ed.com/search.php?search=%v&Submit=Search", fileName))
	if err != nil {
		return "", err
	}
	c.doc = doc

	show, err := c.findShowName()
	if err != nil {
		// Addic7ed did not find the page of the show from the search feature
		results := c.findResults()
		if len(results) == 0 {
			return "", fmt.Errorf("Show not found for filename %v", fileName)
		}
		// If more result, we get the first result
		doc, err := createDocFromURL("http://www.addic7ed.com/" + results[0])
		if err != nil {
			return "", err
		}
		c.doc = doc
		show, err = c.findShowName()
		if err != nil {
			return "", err
		}
	}
	return show, nil
}

// cleanTitle cleans the title of useless words.
// Title are usually of the format "Version BATV, 0.00 MBs", and we want to keep only "BATV"
func cleanTitle(title string) string {
	parts := strings.Split(title, ",")
	clean := title
	if len(parts) >= 0 {
		clean = parts[0]
	}
	parts = strings.Fields(clean)
	if len(parts) >= 2 {
		clean = parts[1]
	} else if len(parts) == 1 {
		clean = parts[0]
	}
	return clean
}

// Parse the string to find words
// The filename is split in words. A word is a a sequence of letters or numbers.
// Every other character is a separator (space, dots, plus, minus...)
func wordsFromString(s string) []string {
	return strings.FieldsFunc(s, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
}

// scoreBestSubVersions give score to subtitles versions
// It searches for similarities in both the filename and the version
// filename and versions are indexed by word, and the more there are common words, the more the version gets a good score
// Similarity is computed from a scoring between word exact matching and word distance (with Jaro/Winkler distance algorithm)
func scoreBestSubVersions(fileName string, subtitlesByVersion map[string]Subtitles) map[string]float64 {
	const weightWhenExactMatch = 5
	wordsFromTitle := wordsFromString(fileName)
	scores := map[string]float64{}

	for version := range subtitlesByVersion {
		versionWords := wordsFromString(version)
		for _, subWordFromTitle := range wordsFromTitle {
			for _, subWordFromVersion := range versionWords {
				var score float64
				if strings.EqualFold(subWordFromVersion, subWordFromTitle) {
					// Give a important weight when we found exactly the same word in the version and the file
					score = weightWhenExactMatch
				}
				// Similarity is a float computed from Jaro/Winkler distance
				// 0 = no similarity at all, 1 = exact same string
				similarityScore := textdistance.JaroWinklerDistance(strings.ToLower(subWordFromVersion), strings.ToLower(subWordFromTitle))
				if val, ok := scores[version]; ok {
					scores[version] = val + score + similarityScore
				} else {
					scores[version] = score + similarityScore
				}
			}
		}
	}
	return scores
}

// findBestSubtitleFromScores returns the best suitable subtitle from the given scores
func findBestSubtitleFromScores(scores map[string]float64, subtitlesByVersion map[string]Subtitles) Subtitle {
	// Get best version from score
	var bestScore float64
	var bestVersion string
	for version, score := range scores {
		if score > bestScore {
			bestVersion = version
			bestScore = score
		}
	}

	// Unable to get the best version from scores, so we get the first to come ¯\_(ツ)_/¯
	if bestVersion == "" {
		// As Go randomizes maps, the "first to come" version may be different between two runs with same input data
		for version := range subtitlesByVersion {
			bestVersion = version
			break
		}
	}

	bestSubs := subtitlesByVersion[bestVersion]

	// From a given version, keep the best subtitle
	// Addic7ed authorizes multiple subtitle of the same version, so we get the most updated one
	var bestSub Subtitle
	for _, sub := range bestSubs {
		if sub.IsUpdated() {
			bestSub = sub
			break
		}
		bestSub = sub
	}
	return bestSub
}

// SearchBest searches in the Addic7ed website for the best suitable subtitle of given episode of a show
// showStr is usually the name of the video file that need to be searched but it could be any search that can be handled by Addic7ed website
// lang is the language of the subtitle
// It returns the episode name and the found subtitle.
func (c *Client) SearchBest(showStr, lang string) (string, Subtitle, error) {
	show, err := c.SearchAll(showStr)
	if err != nil {
		return "", Subtitle{}, err
	}

	subsWithLang := show.Subtitles.Filter(WithLanguage(lang))
	if len(subsWithLang) == 0 {
		return "", Subtitle{}, fmt.Errorf("Unable to find any subtitles for show %q in %q. Check available languages on Addic7ed website and retry", show.Name, lang)
	}

	if len(subsWithLang) == 1 {
		return show.Name, subsWithLang[0], nil
	}

	subsByVersion := subsWithLang.GroupByVersion()

	// Score the different version to find best suitable one
	scores := scoreBestSubVersions(showStr, subsByVersion)

	// From the scores, find the best subtitle possible
	bestSub := findBestSubtitleFromScores(scores, subsByVersion)

	return show.Name, bestSub, nil
}

// SearchAll searches in the Addic7ed website for a given episode of a show
// showStr is usually the name of the video file that need to be searched but it could be any search that can be handled by Addic7ed website
// It returns the episode name and all found subtitles.
func (c *Client) SearchAll(showStr string) (Show, error) {
	showName, err := c.fetchShowPage(showStr)
	if err != nil {
		return Show{}, err
	}
	subtitles := Subtitles{}

	// Search for all HTML table with Addic7ed class tabel95
	c.doc.Find(".tabel95").Each(func(i int, s *goquery.Selection) {
		// Filter only table corresponding to a subtitle version
		if v, ok := s.Attr("align"); ok && v == "center" {
			// Fin the
			title := strings.TrimSpace(s.Find(".NewsTitle").Text())
			s.Find(".language").Each(func(j int, ss *goquery.Selection) {
				language := ss.Text()
				ss.Parent().Find(".buttonDownload").Each(func(k int, sss *goquery.Selection) {
					//fmt.Println("test")
					if val, ok := sss.Attr("href"); ok {
						link := "http://www.addic7ed.com" + val

						version := cleanTitle(title)
						subtitle := Subtitle{
							Version:  version,
							Language: strings.TrimSpace(language),
							Link:     strings.TrimSpace(link),
						}
						subtitles = append(subtitles, subtitle)
					}
				})
			})
		}
	})

	show := Show{
		Name:      showName,
		Subtitles: subtitles,
	}

	return show, nil
}

// Subtitle is a TV-Show subtitle
type Subtitle struct {
	// Language is the Addic7ed language as seen in the website
	Language string
	// Version is the subtitle type/version, usually the name of the teams who ripped the tv show
	Version string
	// Link is the link to the subtitle from Addic7ed website
	Link string
}

func (s Subtitle) String() string {
	return fmt.Sprintf("Link: %v, Version: %v, Language: %v", s.Link, s.Version, s.Language)
}

// IsOriginal checks whether the subtitle is original.
// It means that the subtitle comes with different version and this subtitle is the original one.
func (s Subtitle) IsOriginal() bool {
	return !s.IsUpdated()
}

// Download download the subtitle in-memory, in a closable reader
func (s Subtitle) Download() (io.ReadCloser, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", s.Link, nil)
	if err != nil {
		return nil, err
	}
	// Avoid getting cached pages
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", s.Link) // Without it, the Addic7ed server redirect to the web page instead of dl the srt file

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to reach addic7ed server: %v", err)
	}
	return resp.Body, nil
}

// DownloadTo downloads the subtitle to a given path
func (s Subtitle) DownloadTo(path string) error {
	sub, err := s.Download()
	if err != nil {
		return err
	}
	defer sub.Close()

	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, sub)
	return err
}

// IsUpdated checks whether the subtitle is updated.
// It means that the subtitle comeswith different version and this subtitle is the updated one.
func (s Subtitle) IsUpdated() bool {
	return strings.Contains(s.Link, "updated")
}

// Subtitles is a slice of subtitle
type Subtitles []Subtitle

func (ss Subtitles) String() string {
	subtitles := []string{}
	for _, s := range ss {
		subtitles = append(subtitles, s.String())
	}
	return fmt.Sprintf("[{%v}]", strings.Join(subtitles, "},{"))
}

// Filter filters out subtitles
// To use it, you have to provide a function that returns true for Subtitles to keep, and false to the one to ignore.
// See addic7ed.WithLanguage, addic7ed.WithVersion, addic7ed.WithVersionRegexp for built-in filters
func (ss Subtitles) Filter(filter func(s Subtitle) bool) Subtitles {
	subtitles := Subtitles{}
	for _, subtitle := range ss {
		if filter(subtitle) {
			subtitles = append(subtitles, subtitle)
		}
	}
	return subtitles
}

// GroupBy groups subtitles by a given property from the subtitle
func (ss Subtitles) GroupBy(property func(s Subtitle) string) map[string]Subtitles {
	groupBy := map[string]Subtitles{}
	for _, s := range ss {
		propKey := property(s)
		if val, ok := groupBy[propKey]; ok {
			val = append(val, s)
			groupBy[propKey] = val
		} else {
			groupBy[propKey] = Subtitles{s}
		}
	}
	return groupBy
}

// GroupByVersion groups subtitles by version
func (ss Subtitles) GroupByVersion() map[string]Subtitles {
	return ss.GroupBy(func(s Subtitle) string {
		return s.Version
	})
}

// GroupByLanguage groups subtitles by language
func (ss Subtitles) GroupByLanguage() map[string]Subtitles {
	return ss.GroupBy(func(s Subtitle) string {
		return s.Language
	})
}

// Show defines a TV show with a name and associated subtitle
type Show struct {
	Name      string
	Subtitles Subtitles
}
