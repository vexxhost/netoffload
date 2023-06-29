package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vexxhost/netoffload/pkg/asap2"
)

var VfsCount int

var enableAsap2Cmd = &cobra.Command{
	Use:   "asap2 [device]",
	Short: "Enable network offloading for a device",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := asap2.Enable(args[0], VfsCount)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	enableAsap2Cmd.Flags().IntVarP(&VfsCount, "vfs", "v", 0, "Number of VFs to enable")

	enableCmd.AddCommand(enableAsap2Cmd)
}
