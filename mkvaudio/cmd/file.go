package cmd

import (
	"fmt"
	"os"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/adrianbaraka/goutils/file"
	"github.com/spf13/cobra"
)

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:          "file [path...]",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Short:        "Downmix one or more media files",

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// suggest mkv
		return []string{"mkv"}, cobra.ShellCompDirectiveFilterFileExt
	},

	Run: func(cmd *cobra.Command, args []string) {
		// loop through args and downmix them
		err := false
		for _, path := range args {
			config.Logger.Echof(echo.DefaultColor, echo.Info, "\nProcessing file: '%v'...\n", path)
			downMixed, alreadyStereo, oldFileSize, newFileSize, e := handleFile(path)
			if e != nil {
				config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, e)
				err = true
			} else {
				if downMixed {
					percentReduction := reduction(float64(oldFileSize), float64(newFileSize))
					old := fmt.Sprint(oldFileSize)
					new := fmt.Sprint(newFileSize)

					if !config.bytes {
						old = file.HumanReadableSize(oldFileSize, false)
						new = file.HumanReadableSize(newFileSize, false)
					}
					config.Logger.Echof(echo.Green, echo.Info, "Successfully downmixed file: '%v'. \nOld size: %v, New Size: %v. %.2f%% file size reduction.\n", path, old, new, percentReduction)
				} else {
					if alreadyStereo {
						config.Logger.Echof(echo.Green, echo.Info, "File: '%v' is already stereo. Skipping...\n", path)
					}
				}
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
