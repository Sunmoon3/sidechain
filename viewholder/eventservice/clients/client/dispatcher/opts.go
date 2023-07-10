package dispatcher

import (

	"time"
)

const (
	defaultPeerMonitorPeriod = 600 * time.Second
)


type params struct {
	peerMonitorPeriod    time.Duration

}
// ConnectionReg is a connection registration
type ConnectionReg struct {
	Eventch chan<- *ConnectionEvent
}


func defaultParams( ) *params {
	return &params{peerMonitorPeriod: defaultPeerMonitorPeriod}
}