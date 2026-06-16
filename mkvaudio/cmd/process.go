package cmd

import (
	"github.com/spf13/cobra"
)

// uses app config defined in root.cmd

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process video file(s).",
}

func init() {
	rootCmd.AddCommand(processCmd)

	processCmd.PersistentFlags().BoolVar(&config.dryRun, "dry-run", false, "Show what would be done to the file without any actual modification.")

	processCmd.PersistentFlags().BoolVar(&config.backup, "backup", false, "Do not delete the original file it is kept as a backup.")

	processCmd.PersistentFlags().BoolVar(&config.bytes, "bytes", false, "print sizes in bytes.")
}
