package helpers

import (
	fab "git.labpc.bluarry.top/bluarry/viewholder/eventservice/common"
	"git.labpc.bluarry.top/bluarry/viewholder/internal/fabric/protoutil"
	cb "github.com/hyperledger/fabric-protos-go/common"
	log "github.com/sirupsen/logrus"
)

var logger =log.New()
var AcceptAny fab.BlockFilter = func(block *cb.Block) bool {
	return true
}


func New(headerTypes ...cb.HeaderType) fab.BlockFilter {
	return func(block *cb.Block) bool {
		return hasType(block, headerTypes...)
	}
}

func hasType(block *cb.Block, headerTypes ...cb.HeaderType) bool {
	for i := 0; i < len(block.Data.Data); i++ {
		env, err := protoutil.ExtractEnvelope(block, i)
		if err != nil {
			logger.Errorf("error extracting envelope from block: %s", err)
			continue
		}
		payload, err := protoutil.UnmarshalPayload(env.Payload)
		if err != nil {
			logger.Errorf("error extracting payload from block: %s", err)
			continue
		}
		chdr, err := protoutil.UnmarshalChannelHeader(payload.Header.ChannelHeader)
		if err != nil {
			logger.Errorf("error extracting channel header: %s", err)
			continue
		}
		htype := cb.HeaderType(chdr.Type)
		for _, headerType := range headerTypes {
			if htype == headerType {
				return true
			}
		}
	}
	return false
}