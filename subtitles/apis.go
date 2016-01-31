package subtitles

import (
	"os"
	"strings"

	"github.com/matcornic/subify/subtitles/languages"
	"github.com/matcornic/subify/subtitles/opensubtitles"
	"github.com/matcornic/subify/subtitles/subdb"
	"github.com/olekukonko/tablewriter"
)

// Client defines the interface to get subtitles from API
type Client interface {
	Download(videoPath string, language lang.Language) (subtitlePath string, err error)
	Upload(subtitlePath string, language lang.Language, videoPath string) error
	GetName() string
	GetAliases() []string
}

// Clients is a slice of Client
type Clients []Client

// DefaultAPIs represents the available APIs
// Is also used as the default
var DefaultAPIs = Clients{
	subdb.New(),
	opensubtitles.New(),
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
