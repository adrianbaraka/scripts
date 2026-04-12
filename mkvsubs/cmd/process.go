package cmd

import (
	"github.com/spf13/cobra"
)

// uses app config defined in root.cmd

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process video file(s) to clean or merge subtitles",
}

func init() {
	rootCmd.AddCommand(processCmd)

	processCmd.PersistentFlags().IntVarP(&config.subtitleNumber, "subtitle-number", "s", 1, "The Number Of the subtitle track(can be seen from vlc)")

	processCmd.PersistentFlags().IntVarP(&config.delay, "delay", "d", 0, "Number in milliseconds to delay the chosen subtitle by. Accepts negative numbers. eg -300")

	processCmd.PersistentFlags().BoolVar(&config.dryRun, "dry-run", false, "Show what would be done to the file without any actual modification.")

	processCmd.PersistentFlags().BoolVar(&config.backup, "backup", false, "Do not delete the original file it is kept as a backup.")

	processCmd.PersistentFlags().BoolVar(&config.delImagesubs, "delete-image-subs", false, "Delete image based subs eg vobsub if found the passed subtitle number is image based.")

	processCmd.PersistentFlags().BoolVar(&config.keepOthersubs, "keep-other-subs", false, "If the file has more than one subtitle track preserve them.")

	processCmd.PersistentFlags().BoolVar(&config.mergeScan, "merge-scan", false, "Search for an external subtitle matching the filename. If not found nothing is done to the media file.")

	processCmd.PersistentFlags().StringVar(&config.externalSub, "external-sub", "", "The path to an external sutitle file to merge.")
	processCmd.RegisterFlagCompletionFunc("external-sub", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// suggest srt files
		return []string{"srt"}, cobra.ShellCompDirectiveFilterFileExt
	})
}
