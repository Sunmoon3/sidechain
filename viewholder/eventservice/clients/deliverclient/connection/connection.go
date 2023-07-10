package connection

import (
	"context"
	"fmt"
	"git.labpc.bluarry.top/bluarry/viewholder/chain"
	cb "github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var logger = log.New()

type deliverStream interface {
	grpc.ClientStream
	Send(*cb.Envelope) error
	Recv() (*pb.DeliverResponse, error)
}
// StreamProvider creates a deliver stream
type StreamProvider func(pb.DeliverClient) (stream deliverStream, cancel func(), err error)


var (
	// Deliver creates a Deliver stream
	Deliver = func(client pb.DeliverClient) (deliverStream, func(), error) {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := client.Deliver(ctx)
		return stream, cancel, err
	}

	// DeliverFiltered creates a DeliverFiltered stream
	DeliverFiltered = func(client pb.DeliverClient) (deliverStream, func(), error) {
		ctx, cancel := context.WithCancel(context.Background())
		stream, err := client.DeliverFiltered(ctx)
		return stream, cancel, err
	}
)


type DeliverConnection struct {
	*StreamConnection
	url string
}

func New(config *chain.Config,node *chain.Node,streamProvider StreamProvider) (*DeliverConnection, error){
	logger.Debugf("Connecting to %s...", node.Addr)
	connect, err := NewStreamConnection(config,node,func(grpcconn *grpc.ClientConn) (grpc.ClientStream, func(), error) {
		return streamProvider(pb.NewDeliverClient(grpcconn))
	})
	if err!=nil{
		connect.Close()
		return nil,err
	}

	return &DeliverConnection{
		StreamConnection: connect,
		url: node.Addr,
	},nil
}

func (c *DeliverConnection) deliverStream() deliverStream {
	if c.Stream() == nil {
		return nil
	}
	stream, ok := c.Stream().(deliverStream)
	if !ok {
		panic(fmt.Sprintf("invalid DeliverStream type %T", c.Stream()))
	}
	return stream
}



// Event contains the deliver event as well as the event source
type Event struct {
	SourceURL string
	Event     interface{}
}

// NewEvent returns a deliver event
func NewEvent(event interface{}, sourceURL string) *Event {
	return &Event{
		SourceURL: sourceURL,
		Event:     event,
	}
}

