package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/adrianbaraka/goutils/echo"
)

// uses lang, mkvpropedditexe , logger and runner defined in root.go
func CleanFile(file string) bool {
	tracks, err := getTracks(file)
	if err != nil {
		config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err.Error())
		return false
	}

	lang := fmt.Sprintf("language=%v", config.language)

	var commands [][]string

	// loop through every track (1-indexed, mkvpropedit style) and clear its name
	for i := range tracks {
		commands = append(commands, []string{
			"--edit", fmt.Sprintf("track:%d", i+1),
			"--set", "name=",
		})
	}

	// type-specific tweaks stay as-is
	commands = append(commands,
		[]string{"-e", "track:v1", "-s", "language=und"},
		[]string{"-e", "track:a1", "--set", lang},
		[]string{"-e", "track:s1", "--set", "language=en", "--set", "flag-default=1"},

		[]string{"--delete-attachment", "mime-type:image/png"},
		[]string{"--delete-attachment", "mime-type:image/jpg"},
		[]string{"--delete-attachment", "mime-type:image/jpeg"},

		[]string{"--delete-attachment", "mime-type:application/x-truetype-font"},
		[]string{"--delete-attachment", "mime-type:font/ttf"},
		[]string{"--delete-attachment", "mime-type:application/vnd.ms-opentype"},
		[]string{"--delete-attachment", "mime-type:font/otf"},

		[]string{"--edit", "info", "--set", "title="},
		[]string{"--tags", "all:"},
	)

	finalArgs := []string{file}
	for _, cmdGroup := range commands {
		finalArgs = append(finalArgs, cmdGroup...)
	}

	ok := true
	stdout, err2, code := config.Runner.RunCmd(echo.Debug, config.mkvpropedit.exe, finalArgs...)

	if code < 0 {
		fmt.Println(code)
		ok = false
		config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err2)
	}

	for _, line := range stdout {
		c := echo.Green
		l := echo.Debug
		w := os.Stdout
		if strings.HasPrefix(line, "Warning") {
			c = echo.Yellow
		}
		if strings.HasPrefix(line, "Error") {
			c = echo.Red
			l = echo.Error
			w = os.Stderr
			ok = false
		}
		config.Logger.Fecholn(c, l, w, line)
	}

	return ok
}

func getTracks(file string) ([]mkvTrack, error) {
	stdout, _, code := config.Runner.RunCmd(echo.Debug, config.mkvmerge.exe, "-J", file)
	if code != 0 {
		return nil, fmt.Errorf("mkvmerge -J failed with code %d", code)
	}

	raw := strings.Join(stdout, "\n")
	var info mkvInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return nil, fmt.Errorf("failed to parse mkvmerge JSON: %w", err)
	}
	return info.Tracks, nil
}

type mkvTrack struct {
	Type       string `json:"type"`
	Properties struct {
		UID      uint64 `json:"uid"`
		Language string `json:"language"`
	} `json:"properties"`
}

type mkvInfo struct {
	Tracks []mkvTrack `json:"tracks"`
}

// return OS specific executable name
// also
func getExe(executable string) string {
	if runtime.GOOS == "windows" {
		return fmt.Sprintf("%v.exe", executable)
	}
	return executable
}

func newTool(name, site string) tool {
	return tool{exe: getExe(name), site: site}
}

// Loops through the list if any executable is not found in the system path it is logged and the script exits with a 1 failure.
func verifyTools(tools []tool) {
	notFound := false

	for _, t := range tools {
		path, err := exec.LookPath(t.exe)
		if err != nil {
			config.Logger.Fechof(echo.Red, echo.Error, os.Stderr, "'%v' not found. Check '%v' for installation instructions.\n", t.exe, t.site)
			notFound = true
		} else {
			config.Logger.Echof(echo.Green, echo.Debug, "'%v' found in system path at '%v'.\n", t.exe, path)
		}
	}

	if notFound {
		os.Exit(1)
	}
}
