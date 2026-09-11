/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"mkvaudio/cmd"

	"github.com/spf13/cobra"
)

func main() {
	// ensure all prerun hooks are executed
	cobra.EnableTraverseRunHooks = true
	cmd.Execute()
}
