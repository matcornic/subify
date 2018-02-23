package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.3.0"

// serveCmd represents the serve command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get version of Subify",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
