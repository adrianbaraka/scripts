package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/adrianbaraka/goutils/echo"
)

// uses lang, mkvpropedditexe , logger and runner defined in root.go

func CleanFile(file string) bool {
	// define the commands
	lang := fmt.Sprintf("language=%v", config.language)
	commands := [][]string{
		{"-e", "track:v1", "-s", "name=", "-s", "language=und"},
		{"-e", "track:a1", "-s", "name=", "--set", lang},
		{"-e", "track:s1", "-s", "name=", "--set", "language=en", "--set", "flag-default=1"},

		{"--delete-attachment", "mime-type:image/png"},
		{"--delete-attachment", "mime-type:image/jpg"},
		{"--delete-attachment", "mime-type:image/jpeg"},

		{"--delete-attachment", "mime-type:application/x-truetype-font"},
		{"--delete-attachment", "mime-type:font/ttf"},
		{"--delete-attachment", "mime-type:application/vnd.ms-opentype"},
		{"--delete-attachment", "mime-type:font/otf"},

		{"--edit", "info", "--set", "title="},

		{"--tags", "all:"},
	}

	finalArgs := []string{file}
	for _, cmdGroup := range commands {
		finalArgs = append(finalArgs, cmdGroup...)
	}
	// fmt.Println(finalArgs)
	ok := true
	stdout, err, code := config.Runner.RunCmd(echo.Debug, config.mkvpropeditExe, finalArgs...)

	if code < 0 {
		fmt.Println(code)
		ok = false
		config.Logger.Fecholn(echo.Red, echo.Error, os.Stderr, err)
	}
	//error := false
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
