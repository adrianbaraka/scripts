/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"os/exec"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify the required tools are available in your system path, ie ffmpeg, ffprobe and mkvpropedit.",
	Run: func(cmd *cobra.Command, args []string) {
		verifyTools()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

type tool struct {
	name string
	help string
}

func verifyTools() {
	var tools []tool

	tools = append(tools, tool{config.ffmpegExe, "https://ffmpeg.org/"})
	tools = append(tools, tool{config.ffprobeExe, "https://ffmpeg.org/ffprobe.html"})
	tools = append(tools, tool{config.mkvpropeditExe, "https://mkvtoolnix.download/"})

	notFound := false

	for _, t := range tools {
		path, err := exec.LookPath(t.name)
		if err != nil {
			config.Logger.Fechof(echo.Red, echo.Error, os.Stderr, "'%v' not found. Check '%v' for installation instructions.\n", t.name, t.help)
			notFound = true
		} else {
			config.Logger.Echof(echo.Green, echo.Debug, "'%v' found in system path at '%v'.\n", t.name, path)
		}
	}

	if notFound {
		os.Exit(1)
	}
}
