package cmd

import (
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable network offloading for a device",
}

func init() {
	rootCmd.AddCommand(enableCmd)
}
