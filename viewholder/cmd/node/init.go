package node

import (
	"chainmaker.org/chainmaker/logger/v2"
	"github.com/spf13/cobra"
	"os"
)

var log = logger.GetLogger(config.NODE_CMD)

func InitCmd(_ *cobra.Command, _ []string) {

	ip := config.NODE_DEFAULT_IP
	port := config.NODE_DEFAULT_PORT
	log.Infof("node: %v: %v\n", ip, port)

	log.Info("InitCLTMS...")
	err := sidechainhandle.InitCLTMS(config.CONFIG_PATH)
	if err != nil {
		log.Errorf("err: %v", err)
		os.Exit(1)
	}

	log.Info("InitChainId...")
	err = sidechainhandle.InitChainId(config.CONFIG_PATH, ip+":"+port)
	if err != nil {
		log.Errorf("err: %v", err)
		os.Exit(1)
	}

	log.Info("Init...")
	initFinished := sidechainhandle.InitFinished()
	if !initFinished {
		log.Errorf("侧链初始化InitFinished错误: %v", err)
		os.Exit(1)
	}

	// 5.验证智能合约
	// TODO
	log.Info("开始验证智能合约...")

	// 6.初始化完成
	initedCheck := sidechainhandle.InitedCheck()
	if initedCheck {
		localChainId := sidechainhandle.GetLocalChainId()
		log.Infof("viewholder已经初始化")
		log.Infof("本地链ID: %v", localChainId)
		return
	}
}
