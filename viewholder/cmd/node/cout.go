package node

import (
	"fmt"
	sidechainserver "gitee.com/its_windy_zoe/viewholder_mult/sidechainhandle/server"
	"gitee.com/its_windy_zoe/viewholder_mult/sidechainhandle/server/chaindispatcher"
	"github.com/spf13/cobra"
)

var crossoutCmd = &cobra.Command{
	Use:   "cout",
	Short: "cross out",
	Long:  "cross out",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("args exits")
		}

		err := CrossOutServe(args)
		return err
	},
}

func CrossOutServe(args []string) error {
	dispatcher := chaindispatcher.New()
	sideChainServer := sidechainserver.New(dispatcher)
	errors := sideChainServer.CmdStartOut()
	for {
		select {
		case err := <-errors:
			return err
		}
	}
}
