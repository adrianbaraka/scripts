/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/adrianbaraka/goutils/echo"
	"github.com/spf13/cobra"
)

var (
	verbose   bool
	quiet     bool
	color     string
)

type AppConfig struct {
	Logger    *echo.Logger
	name      string
	season    string
}

var config AppConfig
var allowedColors = []string{"always", "auto", "never"}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename a folder media titles in a consistent clean format",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbosity := echo.Info
		if verbose {
			verbosity = echo.Debug
		}
		if quiet {
			verbosity = echo.Error
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

		config.Logger = echo.NewLogger(verbosity, os.Stdout)
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
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity level to Debug.")
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
