package dispatcher

import (
	"git.labpc.bluarry.top/bluarry/viewholder/chain"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/api"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/common"
	esdispatcher "git.labpc.bluarry.top/bluarry/viewholder/eventservice/service/dispatcher"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var  logger = log.New()



type Dispatcher struct {
	*esdispatcher.Dispatcher
	params
	connection             api.Connection
	connconfig			    *chain.Config
	connectionRegistration *ConnectionReg
	peer 					*chain.Node
	connectionProvider api.ConnectionProvider
	peerMonitorDone        chan struct{}
	lock                   sync.RWMutex
}

func New( cfg *chain.Config,connectionProvider api.ConnectionProvider,opts ...common.Opt)*Dispatcher{
	params := defaultParams()
	common.Apply(params,opts)

	return &Dispatcher{
		connconfig:				cfg,
		Dispatcher:             esdispatcher.New(),
		params:                 *params,
		connectionProvider: connectionProvider,
	}
}

// Start starts the dispatcher
func (ed *Dispatcher) Start() error {
	ed.registerHandlers()

	if err := ed.Dispatcher.Start(); err != nil {
		return errors.WithMessage(err, "error starting client event dispatcher")
	}
	return nil
}


func (ed *Dispatcher) registerHandlers() {
	// Override existing handlers
	ed.RegisterHandler(&esdispatcher.StopEvent{}, ed.HandleStopEvent)

	// Register new handlers
	ed.RegisterHandler(&ConnectEvent{}, ed.HandleConnectEvent)
	ed.RegisterHandler(&DisconnectEvent{}, ed.HandleDisconnectEvent)
	ed.RegisterHandler(&ConnectedEvent{}, ed.HandleConnectedEvent)
	ed.RegisterHandler(&DisconnectedEvent{}, ed.HandleDisconnectedEvent)
	ed.RegisterHandler(&RegisterConnectionEvent{}, ed.HandleRegisterConnectionEvent)
}

// HandleStopEvent handles a Stop event by clearing all registrations
// and stopping the listener
func (ed *Dispatcher) HandleStopEvent(e esdispatcher.Event) {
	// Remove all registrations and close the associated event channels
	// so that the client is notified that the registration has been removed
	ed.clearConnectionRegistration()
	if ed.peerMonitorDone != nil {
		close(ed.peerMonitorDone)
		ed.peerMonitorDone = nil
	}

	ed.Dispatcher.HandleStopEvent(e)
}

//这里是调用连接的地方？
// HandleConnectEvent initiates a connection to the event server
func (ed *Dispatcher) HandleConnectEvent(e esdispatcher.Event) {
	evt := e.(*ConnectEvent)

	if ed.connection != nil {
		// Already connected. No error.
		evt.ErrCh <- nil
		return
	}

	eventch, err := ed.EventCh()
	if err != nil {
		evt.ErrCh <- err
		return
	}
	//获取所有的peers
	//peers, err := ed.discoveryService.GetPeers()
	//if err != nil {
	//	evt.ErrCh <- err
	//	return
	//}

	if len(ed.connconfig.Endorsers) == 0 {
		evt.ErrCh <- errors.New("no peers to connect to")
		return
	}

	//从所有的peer中选择一个合适的peer
	//peer, err := ed.peerResolver.Resolve(peers)
	//if err != nil {
	//	evt.ErrCh <- err
	//	return
	//}
	node:=ed.connconfig.Endorsers[0]

	conn, err := ed.connectionProvider(ed.connconfig,&node)
	if err != nil {
		logger.Warnf("error creating connection: %s", err)
		evt.ErrCh <- errors.WithMessagef(err, "could not create client conn")
		return
	}
	//
	ed.connection = conn
	ed.setConnectedPeer(&node)

	go ed.connection.Receive(eventch)

	evt.ErrCh <- nil
}
func (ed *Dispatcher) setConnectedPeer(peer *chain.Node) {
	ed.lock.Lock()
	defer ed.lock.Unlock()
	ed.peer = peer
}
// HandleRegisterConnectionEvent registers a connection listener
func (ed *Dispatcher) HandleRegisterConnectionEvent(e esdispatcher.Event) {
	evt := e.(*RegisterConnectionEvent)

	if ed.connectionRegistration != nil {
		evt.ErrCh <- errors.New("registration already exists for connection event")
		return
	}

	ed.connectionRegistration = evt.Reg
	evt.RegCh <- evt.Reg
}

// HandleConnectedEvent sends a 'connected' event to any registered listener
func (ed *Dispatcher) HandleConnectedEvent(e esdispatcher.Event) {
	evt := e.(*ConnectedEvent)

	logger.Debugf("Handling connected event: %+v", evt)

	if ed.connectionRegistration != nil && ed.connectionRegistration.Eventch != nil {
		select {
		case ed.connectionRegistration.Eventch <- NewConnectionEvent(true, nil):
		default:
			logger.Warn("Unable to send to connection event channel.")
		}
	}

	if ed.peerMonitorPeriod > 0 {
		ed.peerMonitorDone = make(chan struct{})
		go ed.monitorPeer(ed.peerMonitorDone)
	}
}
func (ed *Dispatcher) monitorPeer(done chan struct{}) {
	logger.Debugf("Starting peer monitor on channel [%s]", "ed.chConfig.ID()")

	ticker := time.NewTicker(ed.peerMonitorPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if ed.disconnected() {
				// Disconnected
				logger.Debugf("Client on channel [%s] has disconnected - stopping disconnect monitor", "ed.chConfig.ID()")
				return
			}
		case <-done:
			logger.Debugf("Stopping block height monitor on channel [%s]", "ed.chConfig.ID()")
			return
		}
	}
}
func (ed *Dispatcher) Connection() api.Connection {
	return ed.connection
}


func (ed *Dispatcher) HandleDisconnectEvent(e esdispatcher.Event) {
	evt := e.(*DisconnectEvent)

	if ed.connection == nil {
		evt.Errch <- errors.New("connection already closed")
		return
	}

	logger.Debug("Closing connection due to disconnect event...")

	ed.connection.Close()
	ed.connection = nil
	//保留待更改连接
	//ed.setConnectedPeer(nil)

	evt.Errch <- nil
}


// HandleDisconnectedEvent sends a 'disconnected' event to any registered listener
func (ed *Dispatcher) HandleDisconnectedEvent(e esdispatcher.Event) {
	evt := e.(*DisconnectedEvent)

	logger.Debugf("Disconnecting from event server: %s", evt.Err)

	if ed.connection != nil {
		ed.connection.Close()
		ed.connection = nil
	}

	if ed.connectionRegistration != nil {
		logger.Debugf("Disconnected from event server: %s", evt.Err)
		select {
		case ed.connectionRegistration.Eventch <- NewConnectionEvent(false, evt.Err):
		default:
			logger.Warn("Unable to send to connection event channel.")
		}
	} else {
		logger.Warnf("Disconnected from event server: %s", evt.Err)
	}

	if ed.peerMonitorDone != nil {
		close(ed.peerMonitorDone)
		ed.peerMonitorDone = nil
	}
}


func (ed *Dispatcher) clearConnectionRegistration() {
	if ed.connectionRegistration != nil {
		logger.Debug("Closing connection registration event channel.")
		close(ed.connectionRegistration.Eventch)
		ed.connectionRegistration = nil
	}
}

/*
待修改
*/
func (ed *Dispatcher) disconnected() bool {
	//connectedPeer := ed.ConnectedPeer()
	//if connectedPeer == nil {
	//	logger.Debugf("Not connected yet")
	//	return false
	//}
	//
	//logger.Debugf("Checking if event client should disconnect from peer [%s] on channel [%s]...", connectedPeer.URL(), ed.chConfig.ID())
	//
	//peers, err := ed.discoveryService.GetPeers()
	//if err != nil {
	//	logger.Warnf("Error calling peer resolver: %s", err)
	//	return false
	//}
	//
	//if !ed.peerResolver.ShouldDisconnect(peers, connectedPeer) {
	//	logger.Debugf("Event client will not disconnect from peer [%s] on channel [%s]...", connectedPeer.URL(), ed.chConfig.ID())
	//	return false
	//}
	//
	//logger.Warnf("The peer resolver determined that the event client should be disconnected from connected peer [%s] on channel [%s]. Disconnecting ...", connectedPeer.URL(), ed.chConfig.ID())
	//
	//if err := ed.disconnect(); err != nil {
	//	logger.Warnf("Error disconnecting event client from peer [%s] on channel [%s]: %s", connectedPeer.URL(), ed.chConfig.ID(), err)
	//	return false
	//}

	logger.Warnf("Successfully disconnected event client from peer [%s] on channel [%s]",ed.peer.Addr,ed.connconfig.Channel)
	return true
}
