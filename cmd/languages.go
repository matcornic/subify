package cmd

import (
	"github.com/matcornic/subify/subtitles/languages"
	"github.com/spf13/cobra"
)

var all bool

// languagesCmd represents the languages command
var languagesCmd = &cobra.Command{
	Use:     "languages",
	Aliases: []string{"lang"},
	Short:   "List available languages",
	Long:    `List available languages`,
	Run: func(cmd *cobra.Command, args []string) {
		lang.Languages.Print(all)
	},
}

func init() {
	listCmd.AddCommand(languagesCmd)
	languagesCmd.PersistentFlags().BoolVar(&all, "all", false, "Shows all languages")
}
