package cmd

import (
	"os"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:  "file [path...]",
	Args: cobra.MinimumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// filter by mkv extension
		return []string{"mkv"}, cobra.ShellCompDirectiveFilterFileExt
	},
	Short: "Clean one or more MKV files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// loop through args and clean them
		err := false
		for _, path := range args {
			config.Logger.Echof(echo.DefaultColor, echo.Info, "Processing file: '%v'...\n", path)
			ok := CleanFile(path)
			if ok {
				config.Logger.Echof(echo.Green, echo.Info, "Successfully cleaned file: '%v'\n", path)
			} else {
				err = true
			}
		}
		if err {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
