package node

import (
	"fmt"
	"github.com/spf13/cobra"
)

var tpsServerCmd = &cobra.Command{
	Use:   "tpserver",
	Short: "Test the Sidechain tps server.",
	Long:  `Start the tps test server that interacts with the sidechain network`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("trailing args detected")
		}
		log.Infof("The Server has been Started")
		//TPS.TPS()
		return nil
	},
}
