package subtitles

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/matcornic/subify/notif"
	"github.com/olekukonko/tablewriter"
	logger "github.com/spf13/jwalterweatherman"
)

// Client defines the interface to get subtitles from API
type Client interface {
	Download(videoPath string, language Language) (subtitlePath string, err error)
	Upload(subtitlePath string, language Language, videoPath string) error
	GetName() string
	GetAliases() []string
}

// Clients is a slice of Client
type Clients []Client

// DefaultAPIs represents the available APIs
// Is also used as the default
var DefaultAPIs = Clients{
	SubDB(),
	OpenSubtitles(),
	Addic7ed(),
}

// InitAPIs sets the order of APIs search from apiAliases
// If alias does not exists, it is not included
func InitAPIs(apiAliases []string) (apis Clients) {
	for _, alias := range apiAliases {
		for _, availableAPI := range DefaultAPIs {
			for _, availableAliases := range availableAPI.GetAliases() {
				if strings.ToLower(strings.TrimSpace(alias)) == strings.ToLower(strings.TrimSpace(availableAliases)) {
					apis = append(apis, availableAPI)
				}
			}
		}
	}
	return
}

// Print prints the clients as nice table
func (c Clients) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Aliases"})
	for _, l := range c {
		values := []string{
			l.GetName(),                        // Name
			strings.Join(l.GetAliases(), ", "), // Aliases
		}
		table.Append(values)
	}
	table.SetAutoWrapText(false)
	table.SetColWidth(50)
	table.SetRowLine(true)
	table.Render() // Send output
}

//String prints a nice representation of clients
func (c Clients) String() (s string) {
	for i, v := range c {
		s = s + v.GetName()
		if (i + 1) < len(c) {
			s = s + ", "
		}
	}
	return
}

// Download the subtitle from the video identified by its path
func Download(videoPath string, apiAliases []string, languages []string) error {
	// APIs to download subtitles.
	var subtitlePath string
	var err error

	// Gets APIs
	a := InitAPIs(apiAliases)
	if len(a) == 0 {
		a = DefaultAPIs
		logger.WARN.Println("No API has been recognized by Subify. Using default:", DefaultAPIs)
	} else if len(apiAliases) != len(a) {
		logger.WARN.Println("Some languages are not recognized. Given:", apiAliases, "Found:", a)
	}

	// Check languages
	l := Languages.GetLanguages(languages)
	if len(l) == 0 {
		logger.ERROR.Println("Languages", languages, "are not available. Pick one ore more from the table below :")
		Languages.Print(false)
		return fmt.Errorf("No languages is available for given languages : %v", languages)
	} else if len(languages) != len(l) {
		logger.WARN.Println("Some languages are not recognized. Given:", languages, "Found:", l.GetDescriptions())
	}

	// Run through languages
browselang:
	for i, lang := range l {
		// Run through different APIs to get the subtitle. Stops when found
		logger.INFO.Println("===> ("+strconv.Itoa(i+1)+") Searching subtitles for", lang.Description, "language")
		for j, api := range a {
			logger.INFO.Println("=> (" + strconv.Itoa(i+1) + "." + strconv.Itoa(j+1) + ") Downloading subtitle with " + api.GetName() + "...")
			subtitlePath, err = api.Download(videoPath, lang)
			if err == nil {
				notif.Info("Subify - I found a subtitle for your video 😎", fmt.Sprintf("Thank you %s ❤️", api.GetName()))
				logger.INFO.Println(lang.Description, "subtitle found and saved to ", subtitlePath)
				break browselang
			} else {
				logger.INFO.Println("Subtitle not found because :", err.Error())
			}
			if (j + 1) < len(a) {
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
		notif.Error("‼️ Subify - I didn't found any subtitle 😭", fmt.Sprintf("No match for your video in : %s. Try later !", a.String()))
		return fmt.Errorf("No %v subtitle found, even after searching in all APIs (%v)", strings.Join(l.GetDescriptions(), ", nor "), a.String())
	}

	return nil
}
