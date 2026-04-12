package cmd

import (
	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:   "folder [name]",
	Short: "Folder containing the media files.",
	Args:  cobra.ExactArgs(1),
	// suggest only one folder
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},

	Run: func(cmd *cobra.Command, args []string) {
		folder := args[0]

		// run the actual function
		runner(folder)
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)

	folderCmd.PersistentFlags().StringVarP(&config.name, "name", "n", "", "The name of the series.")
	folderCmd.PersistentFlags().StringVarP(&config.season, "season", "s", "", "The season eg 01")

	folderCmd.MarkPersistentFlagRequired("name")
	folderCmd.MarkPersistentFlagRequired("season")
}
