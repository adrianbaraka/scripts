package main

import (
	"mkvclean/cmd"

	"github.com/spf13/cobra"
)

func main() {
	// ensure all prerun hooks are executed
	cobra.EnableTraverseRunHooks = true
	cmd.Execute()
}
