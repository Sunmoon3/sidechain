package dispatcher

import (
	"git.labpc.bluarry.top/bluarry/viewholder/chain"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/api"
	clientdsp "git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/client/dispatcher"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/deliverclient/connection"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/common"
	esdispatcher "git.labpc.bluarry.top/bluarry/viewholder/eventservice/service/dispatcher"
	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)


var logger=log.New()

type Dispatcher struct {
	*clientdsp.Dispatcher
}

type dsConnection interface {
	api.Connection
	Send(seekInfo *ab.SeekInfo) error
}


func New(cfg *chain.Config, connectionProvider api.ConnectionProvider,opts ...common.Opt) (*Dispatcher) {
	return &Dispatcher{
		Dispatcher:clientdsp.New(cfg,connectionProvider,opts...),
	}
}


func (ed *Dispatcher) Start() error {
	ed.registerHandlers()
	if err := ed.Dispatcher.Start(); err != nil {
		return errors.WithMessage(err, "error starting deliver event dispatcher")
	}
	return nil
}

func (ed *Dispatcher) registerHandlers() {
	ed.RegisterHandler(&SeekEvent{}, ed.handleSeekEvent)
	ed.RegisterHandler(&connection.Event{}, ed.handleEvent)
}
func (ed *Dispatcher) handleEvent(e esdispatcher.Event) {
	delevent := e.(*connection.Event)
	evt := delevent.Event.(*pb.DeliverResponse)
	switch response := evt.Type.(type) {
	case *pb.DeliverResponse_Status:
		ed.handleDeliverResponseStatus(response)
	case *pb.DeliverResponse_Block:
		ed.HandleBlock(response.Block, delevent.SourceURL)
	case *pb.DeliverResponse_FilteredBlock:
		ed.HandleFilteredBlock(response.FilteredBlock, delevent.SourceURL)
	default:
		logger.Errorf("handler not found for deliver response type %T", response)
	}
}

func (ed *Dispatcher) connection() dsConnection {
	return ed.Dispatcher.Connection().(dsConnection)
}



func (ed *Dispatcher) handleSeekEvent(e esdispatcher.Event) {
	evt := e.(*SeekEvent)

	if ed.Connection() == nil {
		logger.Warn("Unable to register channel since no connection was established.")
		return
	}

	if err := ed.connection().Send(evt.SeekInfo); err != nil {
		evt.ErrCh <- errors.Wrapf(err, "error sending seek info for channel [%s]", "ed.ChannelConfig().ID()")
		//evt.ErrCh <- errors.Wrapf(err, "error sending seek info for channel [%s]", "ed.ChannelConfig().ID()")
	} else {
		evt.ErrCh <- nil
	}
}



func (ed *Dispatcher) handleDeliverResponseStatus(evt *pb.DeliverResponse_Status) {
	logger.Debugf("Got deliver response status event: %#v", evt)

	if evt.Status == cb.Status_SUCCESS {
		return
	}

	logger.Warnf("Got deliver response status event: %#v. Disconnecting...", evt)

	errch := make(chan error, 1)
	ed.Dispatcher.HandleDisconnectEvent(&clientdsp.DisconnectEvent{
		Errch: errch,
	})
	err := <-errch
	if err != nil {
		logger.Warnf("Error disconnecting: %s", err)
	}

	ed.Dispatcher.HandleDisconnectedEvent(disconnectedEventFromStatus(evt.Status))
}

func disconnectedEventFromStatus(status cb.Status) *clientdsp.DisconnectedEvent {
	err := errors.Errorf("got error status from deliver server: %s", status)

	if status == cb.Status_FORBIDDEN {
		return clientdsp.NewFatalDisconnectedEvent(err)
	}
	return clientdsp.NewDisconnectedEvent(err)
}
