package cmd

import (
	"strings"

	"github.com/matcornic/subify/common/utils"
	"github.com/matcornic/subify/subtitles"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var language string
var openVideo bool

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:     "dl <video-path>",
	Aliases: []string{"download"},
	Short:   "Download the subtitles for your video - 'subify dl --help'",
	Long: `Download the subtitles for your video (movie or TV Shows)
Give the path of your video as first parameter and let's go !`,
	Run: func(cmd *cobra.Command, args []string) {
		// Assertions
		utils.VerbosePrintln(logger.INFO,
			"Downloading command called with following parameters : "+strings.Join(args, " "))
		if len(args) != 1 {
			utils.Exit("Video file needed. See usage : 'subify help' or 'subify dl --help'")
		}

		videoPath := args[0]
		subtitles.Download(videoPath)

		if openVideo {
			open.Run(videoPath)
		}
	},
}

func init() {
	dlCmd.Flags().StringVarP(&language, "language", "l", "en", "Language of the subtitle")
	viper.BindPFlag("language", dlCmd.Flags().Lookup("language"))
	dlCmd.Flags().BoolVarP(&openVideo, "open", "o", false,
		"Once the subtitle is donwloaded, open the video with your default video player"+
			` (OSX: "open", Windows: "start", Linux/Other: "xdg-open")`)
	RootCmd.AddCommand(dlCmd)
}
