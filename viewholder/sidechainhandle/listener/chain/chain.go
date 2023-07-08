package chain

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"context"

	"viewholder/config"
	"viewholder/sidechainhandle/events"
)

var log = logger.GetLogger(config.LISTENER_CHAIN)

type ChainListener struct {
	finish chan interface{}
	errch  chan<- error
}

func New(errch chan<- error) *ChainListener {

	return &ChainListener{
		finish: make(chan interface{}),
		errch:  errch,
	}

}

func (cl *ChainListener) Start(evch chan<- interface{}) {
	log.Infof("registChainEvent ...")
	cl.registChainEvents(evch)
}

func (cl *ChainListener) registChainEvents(evch chan<- interface{}) {
	log.Infof("registCrossCCEvent...")

	go cl.registCrossCCEvent(evch)

}

func (cl *ChainListener) registCrossCCEvent(evch chan<- interface{}) {

	chaincodeName := config.USER_CHAINCODE_PATENT
	log.Infof("listen on chaincode %v events %v", chaincodeName, "transferout")

	client, err := examples.CreateChainClientWithSDKConf(config.SDK_CONFIG_PATH)
	if err != nil {
		log.Error("CreateChainClientWithSDKConf:", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	curEvent, err := client.SubscribeContractEvent(ctx, 1, 10, chaincodeName, "")

	if err != nil {
		log.Error(err)
	}

	for {
		select {
		case event, ok := <-curEvent:
			if !ok {
				log.Infof("chan is close!")
				return
			}
			if event == nil {
				log.Error("require not nil")
			}
			contractEventInfo, ok := event.(*common.ContractEventInfo)
			if !ok {
				log.Error("require true")
			}

			log.Infof("send CrossTransferEvent ")
			ev := events.NewCrossTransferEvent(contractEventInfo.TxId, contractEventInfo.EventData)
			evch <- ev
		case <-ctx.Done():
			return
		}

	}

}
