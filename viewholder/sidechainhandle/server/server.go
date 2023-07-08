package server

import (
	"chainmaker.org/chainmaker/logger/v2"
	"viewholder/config"
	"viewholder/sidechainhandle/events"

	"time"
)

var log = logger.GetLogger(config.SIDE_CHAIN_SERVER)

const (
	stopTimeOut = 5 * time.Second
)

type Dispatcher interface {
	StartIn() <-chan error
	StartOut() <-chan error
	EventCh() (chan<- interface{}, error)
}

type SideChainServer struct {
	ds Dispatcher
}

func New(dispatcher Dispatcher) *SideChainServer {
	return &SideChainServer{
		ds: dispatcher,
	}
}

func (s *SideChainServer) CmdStartOut() <-chan error {
	return s.ds.StartOut()
}

func (s *SideChainServer) CmdStartIn() <-chan error {
	return s.ds.StartIn()
}

func (s *SideChainServer) Stop() bool {
	eventCh, err := s.ds.EventCh()
	if err != nil {
		log.Warnf("Error stopping event service: %s", err)
		return false
	}
	errCh := make(chan error, 1)
	eventCh <- events.NewStopEvent(errCh)

	select {
	case err := <-errCh:
		if err != nil {
			log.Warnf("Error while stopping dispatcher: %s", err)
		}
	case <-time.After(stopTimeOut):
		log.Infof("Timed out waiting for dispatcher to stop")
	}
	return true
}
