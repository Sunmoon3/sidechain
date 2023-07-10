package dispatcher

import ab "github.com/hyperledger/fabric-protos-go/orderer"

type SeekEvent struct {
	SeekInfo *ab.SeekInfo
	ErrCh    chan<- error
}

// NewSeekEvent returns a new SeekRequestEvent
func NewSeekEvent(seekInfo *ab.SeekInfo, errch chan<- error) *SeekEvent {
	return &SeekEvent{
		SeekInfo: seekInfo,
		ErrCh:    errch,
	}
}

