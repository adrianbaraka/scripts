package cmd

import (
	"github.com/spf13/cobra"
)

// uses app config defined in root.cmd

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process video file(s).",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		var tools []tool
		tools = append(tools, config.mkvpropedit)
		verifyTools(tools)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)

	// TODO implement dryrun and backup
	processCmd.PersistentFlags().BoolVar(&config.dryRun, "dry-run", false, "Show what would be done to the file without any actual modification. [TODO]")
	processCmd.PersistentFlags().BoolVar(&config.backup, "backup", false, "Do not delete the original file it is kept as a backup. [TODO]")

}
