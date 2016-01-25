package cmd

import "github.com/spf13/cobra"

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List information about something",
	Long:  `List information about something`,
}

func init() {
	RootCmd.AddCommand(listCmd)
}
