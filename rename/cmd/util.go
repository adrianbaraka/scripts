package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrianbaraka/goutils/echo"
)

func fatal(message string) {
	config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, message)
	os.Exit(1)
}

func runner(dir string) {
	if dir == "" {
		fatal("Directory is not defined.")
	}
	if config.season == "" {
		fatal("Season is not defined.")
	}
	if config.name == "" {
		fatal("Series name is not defined.")
	}

	// loop through the dir and dry-run first

	direntry, err := os.ReadDir(dir)
	if err != nil {
		fatal(err.Error())
	}

	var pass = 0
	for {
		for num, file := range direntry {
			ext := filepath.Ext(file.Name())
			oldName := file.Name()
			newName := fmt.Sprintf("%v S%vE%02d%v", config.name, config.season, num+1, ext)

			oldPath := filepath.Join(dir, oldName)
			newPath := filepath.Join(dir, newName)
			if pass == 0 {
				fmt.Printf("Would rename '%v' to '%v'\n", oldPath, newPath)
			} else {
				//fmt.Println("Actual rename")
				err := os.Rename(oldPath, newPath)
				if err != nil {
					config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
				} else {
					config.Logger.Echof(echo.Green, echo.Info, "Successfully renamed '%v' to '%v'\n", oldPath, newPath)
				}
			}
		}

		if pass == 0 {
			// ask for confirmation
			var proceed string
			fmt.Printf("Is this OK: Y [yes] N [no]: ")
			fmt.Scan(&proceed)

			if strings.ToLower(proceed) == "y" || strings.ToLower(proceed) == "yes" {
				fmt.Println("Proceeding with actual renaming...")
			} else {
				fmt.Println("Aborted, no files were renamed.")
				break
			}
		}
		pass++
		if pass > 1 {
			break
		}
	}

}

//for i in {1..9}; do touch "Abbott el S03E0$i.txt"; done
//for i in {10..20}; do touch "Abbott el S03E$i.txt"; done
