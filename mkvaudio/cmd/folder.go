package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/adrianbaraka/goutils/file"
	"github.com/spf13/cobra"
)

// folderCmd represents the folder command
var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "The folder to search and downmix media files.",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveFilterDirs
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	Run: func(cmd *cobra.Command, args []string) {
		dir := args[0]

		// Thread-safe atomic counters
		var count atomic.Int64
		var downmixedCount atomic.Int64
		var oldsize atomic.Int64
		var newsize atomic.Int64
		var failed atomic.Int64
		var wouldDownmix atomic.Int64

		// Limit concurrency to half the available CPU cores (minimum 1)
		maxWorkers := max(runtime.NumCPU()/2, 1)
		semaphore := make(chan struct{}, maxWorkers)
		var wg sync.WaitGroup

		// Recursively walk through all files and subdirectories
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				// Log the error accessing the path and continue walking other paths
				config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
				return nil
			}

			config.Logger.Echof(echo.DefaultColor, echo.Trace, "Checking: '%v'\n", path)
			if d.IsDir() || filepath.Ext(path) != ".mkv" {
				config.Logger.Echoln(echo.DefaultColor, echo.Trace, "Skipping ", path)
				return nil
			}

			// Block if semaphore is full
			semaphore <- struct{}{}
			wg.Add(1)

			go func(fname string) {
				defer wg.Done()
				defer func() { <-semaphore }() // Release slot

				config.Logger.Echof(echo.DefaultColor, echo.Info, "\nProcessing file: '%v'...\n", fname)
				count.Add(1)

				downMixed, alreadyStereo, oldFileSize, newFileSize, err := handleFile(fname)
				if err != nil {
					failed.Add(1)
					config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr,fname, ":", err)
					return
				}

				oldsize.Add(oldFileSize)
				newsize.Add(newFileSize)

				if downMixed {
					downmixedCount.Add(1)
					config.Logger.Echof(echo.Green, echo.Info, "Successfully downmixed file: '%v'\n", fname)
				} else {
					if alreadyStereo {
						config.Logger.Echof(echo.Green, echo.Info, "File: '%v' is already stereo. Skipping...\n", fname)
					} else {
						wouldDownmix.Add(1)
					}
				}
			}(path)

			return nil
		})

		if err != nil {
			config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
			os.Exit(1)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Read final atomic values for reporting
		finalCount := count.Load()
		finalFailed := failed.Load()
		finalDownmixed := downmixedCount.Load()
		finalOldSize := oldsize.Load()
		finalNewSize := newsize.Load()
		finalWouldDownMix := wouldDownmix.Load()

		percentReduction := reduction(float64(finalOldSize), float64(finalNewSize))
		old := fmt.Sprint(finalOldSize)
		new := fmt.Sprint(finalNewSize)

		if !config.bytes {
			old = file.HumanReadableSize(finalOldSize, false)
			new = file.HumanReadableSize(finalNewSize, false)
		}

		config.Logger.Echof(echo.DefaultColor, echo.Info, "------------%v mkv files found and processed %v failed.------------\n", finalCount, finalFailed)
		if config.dryRun {
			// put error just so it shows up even with -q option
			config.Logger.Echof(echo.DefaultColor, echo.Error, "Would Downmix %v files.\n", finalWouldDownMix)
		} else {
			config.Logger.Echof(echo.DefaultColor, echo.Info, "%v files downmixed. Old size: %v, New Size: %v. %.2f%% file size reduction.\n", finalDownmixed, old, new, percentReduction)
		}
		if finalFailed > 0 {
			os.Exit(1)
		}
	},
}

func init() {
	processCmd.AddCommand(folderCmd)
}
