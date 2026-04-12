package cmd

import (
	"os"
	"path/filepath"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "The folder to search and clean mkv files.",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]
		count := 0
		failed := 0

		direntry, err := os.ReadDir(dir)
		if err != nil {
			config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
			os.Exit(1)
		}
		// loop through all files in dir
		for _, entry := range direntry {
			f := entry.Name()

			filename := filepath.Join(dir, f)

			config.Logger.Echof(echo.DefaultColor, echo.Trace, "Checking: '%v'\n", filename)
			if entry.IsDir() || filepath.Ext(filename) != ".mkv" {
				config.Logger.Echoln(echo.DefaultColor, echo.Trace, "Skipping ", filename)
				continue
			}
			config.Logger.Echof(echo.DefaultColor, echo.Info, "Processing file: '%v'...\n", filename)
			count++
			err := handleFile(filename)
			if err != nil {
				failed++
				config.Logger.Echoln(echo.Red, echo.Error, err)
			} else {
				config.Logger.Echof(echo.Green, echo.Info, "Successfully processed file: '%v'\n", filename)
			}
		}
		config.Logger.Echof(echo.DefaultColor, echo.Info, "------------%v mkv files found and processed %v failed.------------\n", count, failed)
		if failed > 0 {
			os.Exit(1)
		}

	},
}

func init() {
	processCmd.AddCommand(folderCmd)
}
