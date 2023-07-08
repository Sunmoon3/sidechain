package eventprocessor

import (
	"github.com/pkg/errors"
	"reflect"
	"sync/atomic"
	"viewholder/config"
	"viewholder/sidechainhandle/events"
)

const (
	eventConsumerBufferSize = 100
)

const (
	dispatcherStateInitial = iota
	dispatcherStateStarted
	dispatcherStateStopped
)

var log = logger.GetLogger(config.EVENT_PROCESSOR)

type Handler func(events.Event)

type EventProcessor struct {
	state    int32
	eventch  chan interface{}
	handlers map[reflect.Type]Handler
	errch    chan error
}

func New() *EventProcessor {
	return &EventProcessor{
		eventch:  make(chan interface{}, eventConsumerBufferSize),
		errch:    make(chan error, eventConsumerBufferSize),
		handlers: make(map[reflect.Type]Handler),
		state:    dispatcherStateInitial,
	}
}

func (ep *EventProcessor) RegisterHandler(t interface{}, h Handler) {
	htype := reflect.TypeOf(t)
	if _, ok := ep.handlers[htype]; !ok {
		log.Debugf("Registering handler for %s on EventProcessor %T", htype, ep)
		ep.handlers[htype] = h
	} else {
		log.Debugf("Cannot register handler %s on EventProcessor %T since it's already registered", htype, ep)
	}
}

func (ep *EventProcessor) EventCh() (chan<- interface{}, error) {
	state := ep.getState()
	if state == dispatcherStateStarted {
		return ep.eventch, nil
	}
	return nil, errors.Errorf("dispatcher not started - Current state [%d]", state)
}

func (ep *EventProcessor) ErrorCh() chan error {
	return ep.errch
}

func (ep *EventProcessor) setState(expectedState, newState int32) bool {
	return atomic.CompareAndSwapInt32(&ep.state, expectedState, newState)
}

func (ep *EventProcessor) getState() int32 {
	return atomic.LoadInt32(&ep.state)
}

func (ep *EventProcessor) Start() error {
	if !ep.setState(dispatcherStateInitial, dispatcherStateStarted) {
		return errors.New("cannot start dispatcher since it's not in its initial state")
	}

	go func() {
		for {
			if ep.getState() == dispatcherStateStopped {
				break
			}
			log.Debug("Listening for events...")
			e, ok := <-ep.eventch
			if !ok {
				break
			}

			log.Debugf("Received event: %+v", reflect.TypeOf(e))

			if handler, ok := ep.handlers[reflect.TypeOf(e)]; ok {
				log.Debugf("Dispatching event: %+v", reflect.TypeOf(e))
				handler(e)
			} else {
				log.Errorf("Handler not found for: %s", reflect.TypeOf(e))
			}
		}
		log.Debug("Exiting event dispatcher")
	}()
	return nil
}

func (ep *EventProcessor) Stop(errch chan<- error) {
	log.Debugf("Stopping dispatcher...")
	if !ep.setState(dispatcherStateStarted, dispatcherStateStopped) {
		log.Warn("Cannot stop event dispatcher since it's already stopped.")
		errch <- errors.New("dispatcher already stopped")
		return
	}
}
