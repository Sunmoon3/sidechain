package deliverclient

import (
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/api"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/deliverclient/seek"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/common"
	"time"
)

type params struct{
	seekType     seek.Type
	fromBlock    uint64
	respTimeout  time.Duration
	connProvider api.ConnectionProvider
}


func defaultParams() *params {
	return &params{
		connProvider: deliverProvider,
		respTimeout:  5 * time.Second,
	}
}


// WithSeekType specifies the point from which block events are to be received.
func WithSeekType(value seek.Type) common.Opt {
	return func(p common.Params) {
		if setter, ok := p.(seekTypeSetter); ok {
			setter.SetSeekType(value)
		}
	}
}

type seekTypeSetter interface {
	SetSeekType(value seek.Type)
}
