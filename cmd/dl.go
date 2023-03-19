package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/matcornic/subify/common/utils"
	"github.com/matcornic/subify/subtitles"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	logger "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var openVideo bool

var notify bool

// dlCmd represents the dl command
var dlCmd = &cobra.Command{
	Use:     "dl <video-path>",
	Aliases: []string{"download"},
	Short:   "Download the subtitles for your video - 'subify dl --help'",
	Long: `Download the subtitles for your video (movie or TV Shows)
Give the path of your video as first parameter and let's go !`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			utils.Exit("At least one video file needed. See usage : 'subify help' or 'subify dl --help'")
		}

		targets := enumerateTargets(args)

		for _, videoPath := range targets {
			utils.VerbosePrintln(logger.INFO, "Given video file is "+videoPath)

			apis := strings.Split(viper.GetString("download.apis"), ",")
			languages := strings.Split(viper.GetString("download.languages"), ",")
			err := subtitles.Download(videoPath, apis, languages, notify)
			if err != nil {
				utils.VerbosePrintln(logger.ERROR, "Sadly, we could not download any subtitle for \""+videoPath+"\". Try another time or contribute to the apis. See 'subify upload -h'")
			}
			if openVideo {
				err := open.Run(videoPath)
				if err != nil {
					utils.ExitPrintError(err, "Sadly, we could not open video: %s", videoPath)
				}
			}
		}
	},
}

// enumerateTargets checks if there are any directories in the argument list:
// if there are, it recursively adds subfiles that are videos
func enumerateTargets(paths []string) []string {
	res := []string{}

	for _, path := range paths {
		// Check if the file is a directory. If it is, use all videos inside it as targets
		stat, err := os.Stat(path)
		if err != nil {
			utils.VerbosePrintln(logger.ERROR, err.Error())
			continue
		}

		if stat.IsDir() {
			res = append(res, getTargetsFromDirectory(path)...)
		} else {
			res = append(res, path)
		}
	}

	return res
}

// getTargetsFromDirectory receives a directory path and recursively returns
// all it's video files and subfiles
func getTargetsFromDirectory(path string) []string {
	res := []string{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		utils.VerbosePrintln(logger.ERROR, err.Error())
		return nil
	}

	for _, file := range files {
		fullPath := filepath.Join(path, file.Name())

		if !file.IsDir() && isVideoFilename(file.Name()) {
			res = append(res, fullPath)
		} else if file.IsDir() {
			res = append(res, getTargetsFromDirectory(fullPath)...)
		}
	}

	return res
}

// isVideoFilename checks if a filename ends in a video extension
func isVideoFilename(fname string) bool {
	ext := filepath.Ext(fname)
	for _, v := range viper.GetStringSlice("download.extensions") {
		if "."+v == ext {
			return true
		}
	}
	return false
}

func init() {
	dlCmd.Flags().StringP("languages", "l", "en", "Languages of the subtitle separate by a comma (First to match is downloaded). Available languages at 'subify list languages'")
	dlCmd.Flags().StringP("apis", "a", "SubDB,OpenSubtitles,Addic7ed", "Overwrite default searching APIs behavior, hence the subtitles are downloaded. Available APIs at 'subify list apis'")
	dlCmd.Flags().StringSliceP("extensions", "e", []string{"wmv", "mov", "webm", "mkv", "avi", "mp4"}, "List of file extensions that should be considered videos.")
	dlCmd.Flags().BoolVarP(&openVideo, "open", "o", false,
		"Once the subtitle is downloaded, open the video with your default video player"+
			` (OSX: "open", Windows: "start", Linux/Other: "xdg-open")`)
	dlCmd.Flags().BoolVarP(&notify, "notify", "n", true, "Display desktop notification")
	_ = viper.BindPFlag("download.languages", dlCmd.Flags().Lookup("languages"))
	_ = viper.BindPFlag("download.apis", dlCmd.Flags().Lookup("apis"))
	_ = viper.BindPFlag("download.extensions", dlCmd.Flags().Lookup("extensions"))

	RootCmd.AddCommand(dlCmd)
}
