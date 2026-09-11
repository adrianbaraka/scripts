/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Verify the required tools are available in your system path, ie ffmpeg, ffprobe and mkvpropedit.",
	Run: func(cmd *cobra.Command, args []string) {
		var tools []tool
		tools = append(tools, config.mkvextract)
		tools = append(tools, config.mkvmerge)
		verifyTools(tools)
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
