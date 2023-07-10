package client

import (
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/client/dispatcher"
	"time"
)

type params struct {
	connEventCh             chan *dispatcher.ConnectionEvent
	reconnInitialDelay      time.Duration
	timeBetweenConnAttempts time.Duration
	respTimeout             time.Duration
	eventConsumerBufferSize uint
	maxConnAttempts         uint
	maxReconnAttempts       uint
	permitBlockEvents       bool
	reconn                  bool
}


func defaultParams() *params {
	return &params{
		eventConsumerBufferSize: 100,
		reconn:                  true,
		maxConnAttempts:         1,
		maxReconnAttempts:       0, // Try forever
		reconnInitialDelay:      0,
		timeBetweenConnAttempts: 5 * time.Second,
		respTimeout:             5 * time.Second,
	}
}
