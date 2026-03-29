package cmd

import (
	"fmt"
	"os"
	"runtime"
	"slices"

	"github.com/adrianbaraka/goutils/cli"
	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

var (
	verbose        bool
	quiet          bool
	mkvpropeditexe string
	language       string
	color          string
)

type AppConfig struct {
	Logger         *echo.Logger
	Runner         *cli.RunCmdConfig
	mkvpropeditExe string
	language       string
}

var config AppConfig

var allowedColors = []string{"always", "auto", "never"}

func mkvpath() string {
	if runtime.GOOS == "windows" {
		return "mkvpropedit.exe"
	}
	return "mkvpropedit"
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mkvclean",
	Short: "Strip fonts, images, and track names from MKV files",
	Long: `mkvclean is a utility to sanitize Matroska files using mkvpropedit. 
It recursively (or individually) processes files to remove:
  - All embedded fonts and images (attachments)
  - Metadata titles from Video, Audio, and Subtitle tracks
  - Global tags and statistics

Example:
  mkvclean file movie.mkv
  mkvclean folder ./movies --color never`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//TODO conf the help files
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// configure options

		verbosity := echo.Info
		if verbose {
			verbosity = echo.Debug
		}
		if quiet {
			verbosity = echo.Error
		}

		// handle color
		if !slices.Contains(allowedColors, color) {
			return fmt.Errorf("Invalid --color value: %s (choose from %v)", color, allowedColors)
		}
		// these env vars are only active for the duration of the script
		if color == "always" {
			os.Setenv("FORCE_COLOR", "true")
		}
		if color == "never" {
			os.Unsetenv("FORCE_COLOR")
			os.Setenv("NO_COLOR", "true")
		}

		config.Logger = echo.NewLogger(verbosity, os.Stdout)
		config.Runner = cli.NewRunner(verbosity, false, true, false)

		config.mkvpropeditExe = mkvpropeditexe
		config.language = language

		return nil
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
	rootCmd.Version = "0.2.2"

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity level to Debug.")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Decrease verbosity level to only errors.")
	// mark verbose and quiet as mutually exclusive
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "quiet")

	rootCmd.PersistentFlags().StringVarP(&mkvpropeditexe, "mkvpropedit", "m", mkvpath(), "Path to the mkvpropedit executable if it is not in the $PATH")

	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "en", "The language of the first audio track. See https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes")

	// color flag
	rootCmd.PersistentFlags().StringVar(&color, "color", "auto", "Colorize output (always, auto, never)")
	// tab completion for color
	rootCmd.RegisterFlagCompletionFunc("color", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allowedColors, cobra.ShellCompDirectiveNoFileComp
	})
}

// TODO  tests
