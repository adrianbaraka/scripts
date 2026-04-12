package cmd

import (
	"os"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:          "file [path...]",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Clean one or more MKV files",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// filter by mkv extension
		return []string{"mkv"}, cobra.ShellCompDirectiveFilterFileExt
	},

	Run: func(cmd *cobra.Command, args []string) {
		// loop through args and clean them
		err := false
		for _, path := range args {
			config.Logger.Echof(echo.DefaultColor, echo.Info, "\nProcessing file: '%v'...\n", path)
			e := handleFile(path)
			if e == nil {
				config.Logger.Echof(echo.Green, echo.Info, "Successfully processed file: '%v'\n", path)
			} else {
				config.Logger.Echoln(echo.Red, echo.Debug, e)
				err = true
			}
		}
		if err {
			os.Exit(1)
		}
	},
}

func init() {
	processCmd.AddCommand(fileCmd)
}
