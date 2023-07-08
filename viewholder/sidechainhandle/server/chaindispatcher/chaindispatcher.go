package chaindispatcher

import (
	"viewholder/config"
	"viewholder/sidechainhandle/cltmshandlesend"
	"viewholder/sidechainhandle/events"
	"viewholder/sidechainhandle/listener/chain"
	"viewholder/sidechainhandle/server/crossdispatcher"
)

var log = logger.GetLogger(config.CHAIN_DISPATCHER)

type ChainDispatcher struct {
	*crossdispatcher.CrossDispatcher
	chainlis *chain.ChainListener
}

func New() *ChainDispatcher {
	return &ChainDispatcher{
		CrossDispatcher: crossdispatcher.New(),
	}
}

func (chd *ChainDispatcher) StartIn() <-chan error {
	errorCh := chd.ErrorCh()

	log.Debugf("The Server's EventProcessor is starting...")
	err := chd.EventProcessor.Start()
	if err != nil {
		errorCh <- err
		return errorCh
	}
	log.Debugf("The Server's EventProcessor is started")

	chd.CrossDispatcher.StartLanguageReceived()
	return errorCh
}

func (chd *ChainDispatcher) StartOut() <-chan error {

	errorCh := chd.ErrorCh()

	chd.registeSendhandlers()

	log.Debugf("The Server's EventProcessor is starting...")
	err := chd.EventProcessor.Start()
	if err != nil {
		errorCh <- err
		return errorCh
	}
	log.Debugf("The Server's EventProcessor is started")

	log.Infof("The Server's Chainlistener is starting...")
	chd.chainlis = chain.New(errorCh)

	eventCh, err := chd.EventCh()

	if err != nil {
		errorCh <- err
		return errorCh
	}
	chd.chainlis.Start(eventCh)
	log.Infof("The Server's Chainlistener is started")

	chd.CrossDispatcher.StartLanguageGenerated()

	return errorCh
}

func (chd *ChainDispatcher) registeSendhandlers() {
	chd.RegisterHandler(&events.CrossTransferEvent{}, chd.handleCrossTransferEvent)
	chd.RegisterHandler(&events.StopEvent{}, chd.HandleStopEvent)
}

func (chd *ChainDispatcher) handleCrossTransferEvent(e events.Event) {
	cte := e.(*events.CrossTransferEvent)

	errch := chd.ErrorCh()
	eventch, err := chd.EventCh()
	if err != nil {
		errch <- err
	}
	log.Infof("receive CrossTransferEvent,and start Solving...")

	cltmshandlesend.CrossOutAggre(cte, eventch, errch)
	log.Infof("CrossTransferEvent solve success")
}

func (chd *ChainDispatcher) HandleStopEvent(e events.Event) {
	event := e.(*events.StopEvent)

	chd.chainlis.Stop()

	chd.CrossDispatcher.Stop(event.ErrCh)

	event.ErrCh <- nil
}
