package dispatcher

import (
	"time"
)

type params struct {
	eventConsumerBufferSize           uint
	eventConsumerTimeout              time.Duration
	initialLastBlockNum               uint64
	initialBlockRegistrations         []*BlockReg
	initialFilteredBlockRegistrations []*FilteredBlockReg
	initialCCRegistrations            []*ChaincodeReg
	initialTxStatusRegistrations      []*TxStatusReg
}


func defaultParams() *params {
	return &params{
		eventConsumerBufferSize: 100,
		eventConsumerTimeout:    500 * time.Millisecond,
	}
}

