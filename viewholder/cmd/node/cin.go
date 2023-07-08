package node

import (
	"fmt"
	"github.com/spf13/cobra"
)

var crossinCmd = &cobra.Command{
	Use:   "cin",
	Short: "cross in",
	Long:  "cross in",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("args exits")
		}

		err := CrossInServe(args)
		return err
	},
}

func CrossInServe(args []string) error {
	dispatcher := chaindispatcher.New()
	sideChainServer := sidechainserver.New(dispatcher)
	errors := sideChainServer.CmdStartIn()
	for {
		select {
		case err := <-errors:
			return err
		}
	}
}
