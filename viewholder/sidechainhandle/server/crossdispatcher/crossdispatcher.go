package crossdispatcher

import (
	"fmt"
	"viewholder/config"
	"viewholder/sidechainhandle/events"
	"viewholder/sidechainhandle/listener/cross"
	"viewholder/sidechainhandle/server/eventprocessor"
)

var log = logger.GetLogger(config.CROSS_DISPATCHER)

type CrossDispatcher struct {
	*eventprocessor.EventProcessor                     
	crslis                         *cross.CrossListener 
}

func New() *CrossDispatcher {
	c := &CrossDispatcher{
		EventProcessor: eventprocessor.New(),
	}
	errch := c.ErrorCh()
	c.crslis = cross.NewCin(errch)

	return c
}

func (cd *CrossDispatcher) StartLanguageGenerated() {
	cd.registerLanguageGeneratedhandlers()
	errorCh := cd.ErrorCh()
	eventCh, err := cd.EventCh()
	if err != nil {
		errorCh <- fmt.Errorf("error when to get eventch ，error:%v", err)
		return
	}
	log.Infof("The Server's Crosslistener is starting...")
	cd.crslis.Start(eventCh)
	log.Infof("The Server's Crosslistener is started")
}

func (cd *CrossDispatcher) StartLanguageReceived() {
	cd.registerLanguageReceivedhandlers()
}

func (cd *CrossDispatcher) Stop(errch chan<- error) {
	cd.crslis.Stop()
	cd.EventProcessor.Stop(errch)
}

func (cd *CrossDispatcher) registerLanguageGeneratedhandlers() {
	log.Info("Debug! Gen")
	cd.RegisterHandler(&events.CrosslanguageGeneratedEvent{}, cd.handleCrosslgGenedEvent)
}

func (cd *CrossDispatcher) handleCrosslgGenedEvent(e events.Event) {
	clge := e.(*events.CrosslanguageGeneratedEvent)
	go cd.crslis.SendtoSC(clge.Req)
}

func (cd *CrossDispatcher) registerLanguageReceivedhandlers() {
	log.Info("Debug! Rec")
	cd.RegisterHandler(&events.CrossLanguageRecvedEvent{}, cd.handleCrosslgRecvedEvent)
}

func (cd *CrossDispatcher) handleCrosslgRecvedEvent(event events.Event) {
	clrcev := event.(*events.CrossLanguageRecvedEvent)

	log.Infof("start crossIn...")

	cltmshandlerec.CrossIn(clrcev, cd.ErrorCh())
}
