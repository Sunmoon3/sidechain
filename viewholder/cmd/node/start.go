package node

import (
	"fmt"
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	nodeStartCmd.AddCommand(crossoutCmd)
	nodeStartCmd.AddCommand(crossinCmd)
	nodeStartCmd.AddCommand(tpsServerCmd)
	return nodeStartCmd
}

var nodeStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the node",
	Long:  "Start the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("args exits")
		}
		log.Infof("viewholder success!")
		return nil
	},
}
