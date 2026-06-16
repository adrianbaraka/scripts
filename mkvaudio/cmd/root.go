package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/adrianbaraka/goutils/cli"
	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

var (
	verbosityCount int
	quiet          bool
	color          string
)

type AppConfig struct {
	Logger *echo.Logger
	Runner *cli.RunCmdConfig
	backup bool
	dryRun bool
	bytes  bool

	// executables
	ffmpegExe      string
	ffprobeExe     string
	mkvpropeditExe string
}

var config AppConfig
var allowedColors = []string{"always", "auto", "never"}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mkvaudio",
	Short: "Downmix audio files to stereo.",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		verbosity := echo.Info

		if quiet {
			verbosity = echo.Error
		} else {
			// Switch based on how many 'v's were passed
			switch verbosityCount {
			case 0:
				verbosity = echo.Info
			case 1:
				verbosity = echo.Debug // -v
			default:
				verbosity = echo.Trace // -vv
			}
		}

		// handle color
		if !slices.Contains(allowedColors, color) {
			fmt.Fprintf(os.Stderr, "Invalid --color value: %s (choose from %v)\n", color, allowedColors)
			os.Exit(1)
		}
		// these env vars are only active for the duration of the script
		if color == "always" {
			os.Setenv("FORCE_COLOR", "true")
		}
		if color == "never" {
			os.Unsetenv("FORCE_COLOR")
			os.Setenv("NO_COLOR", "true")
		}

		config.ffmpegExe = getExe("ffmpeg")
		config.ffprobeExe = getExe("ffprobe")
		config.mkvpropeditExe = getExe("mkvpropedit")

		config.Logger = echo.NewLogger(verbosity, os.Stdout)
		// the mktoolinix clis dont write errors to stderr so can't stream output in color so either color or stream but not both
		// runner is used when capturing output is needed streamer to stream only

		config.Runner = cli.NewRunner(verbosity, true, true, false)
		//config.Streamer = cli.NewRunner(verbosity, true, false, true)
		//config.Streamer = config.Runner.RunCmdStreamer()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbosityCount, "verbose", "v", "Increase verbosity level (-v for Debug, -vv for Trace).")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Decrease verbosity level to only errors.")
	// mark verbose and quiet as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "quiet")

	// color flag
	rootCmd.PersistentFlags().StringVar(&color, "color", "auto", "Colorize output (always, auto, never)")
	// tab completion for color
	rootCmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allowedColors, cobra.ShellCompDirectiveNoFileComp
	})

}
