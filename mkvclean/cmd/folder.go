package cmd

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:  "folder",
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Short: "The folder to recursively search and clean mkv files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		directory := args[0]
		count := 0
		failed := 0
		err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			config.Logger.Echof(echo.DefaultColor, echo.Debug, "Checking: '%v'\n", path)
			if !d.IsDir() && filepath.Ext(path) == ".mkv" {

				config.Logger.Echof(echo.DefaultColor, echo.Info, "Found file: '%v'. Cleaning...\n", path)
				ok := CleanFile(path)
				if ok {
					count++
					config.Logger.Echof(echo.Green, echo.Info, "Successfully cleaned file: '%v'\n", path)
				} else {
					failed++
				}
				//run.CleanFile(path, mkvpropeditexe, language, *runner, logger)
			}
			return nil

		})
		if err != nil {
			config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
			os.Exit(1)
		} else {
			config.Logger.Echof(echo.DefaultColor, echo.Info, "------------%v mkv files found and cleaned %v failed.------------\n", count, failed)
			if failed > 0 {
				os.Exit(1)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// folderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// folderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
