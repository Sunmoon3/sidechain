package connection

import (
	"errors"
	"git.labpc.bluarry.top/bluarry/viewholder/chain"
	clientdsp "git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/client/dispatcher"
	"git.labpc.bluarry.top/bluarry/viewholder/internal/fabric/protoutil"
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-protos-go/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"google.golang.org/grpc"
	"io"
	"sync/atomic"
)
type StreamProviderx func(conn *grpc.ClientConn) (grpc.ClientStream, func(), error)

type StreamConnection struct {
	stream   	grpc.ClientStream
	conn        *grpc.ClientConn
	node		*chain.Node
	config 		*chain.Config
	cancel   	func()
	done		int32
}
func NewStreamConnection(config *chain.Config,node *chain.Node,streamProvider StreamProviderx)(*StreamConnection, error){
	conn, err := chain.DailConnection(*node, logger)
	if err!=nil{
		return nil,err
	}

	stream, cancel, err := streamProvider(conn)
	if err!=nil{
		cancel()
		return nil, err
	}
	if stream == nil {
		return nil, errors.New("unexpected nil stream received from provider")
	}


	return &StreamConnection{
		stream: stream,
		cancel: cancel,
		conn:   conn,
		node:   node,
		config:	config,
	},nil
}


func (c *StreamConnection) Close() {
	if !c.setClosed() {
		logger.Debug("Already closed")
		return
	}

	logger.Debug("Releasing connection....")
	err:=c.conn.Close()
	if err!=nil{
		logger.Errorf("close grpc error: %v",err)
	}
	logger.Debug("... connection successfully closed.")
}


func (c *StreamConnection) setClosed() bool {
	return atomic.CompareAndSwapInt32(&c.done, 0, 1)
}
// Closed returns true if the connection has been closed
func (c *StreamConnection) Closed() bool {
	return atomic.LoadInt32(&c.done) == 1
}


func (c *StreamConnection) Stream() grpc.ClientStream {
	return c.stream
}
// Send sends a seek request to the deliver server
func (c *DeliverConnection) Send(seekInfo *ab.SeekInfo) error {
	if c.Closed() {
		return errors.New("connection is closed")
	}

	logger.Debugf("Sending %#v", seekInfo)

	env, err := c.createSignedEnvelope(seekInfo)
	if err != nil {
		return err
	}

	return c.deliverStream().Send(env)
}
func (c *DeliverConnection) createSignedEnvelope(msg proto.Message) (*cb.Envelope, error) {
	// TODO: Do we need to make these configurable?
	var msgVersion int32
	var epoch uint64

	payloadChannelHeader := protoutil.MakeChannelHeader(cb.HeaderType_DELIVER_SEEK_INFO, msgVersion,c.config.Channel, epoch)


	//先设置为空试试
	//来自 fab.EndpointConfig的TLSClientCerts方法去certs hash
	//payloadChannelHeader.TlsCertHash = c.TLSCertHash()

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	//identity, err := c.Context().Serialize()
	//if err != nil {
	//	return nil, err
	//}

	//nonce, err := crypto.GetRandomNonce()
	//if err != nil {
	//	return nil, err
	//}

	cry,err:=c.config.LoadCrypto()
	if err != nil {
		return nil, err
	}
	payloadSignatureHeader,err:=cry.NewSignatureHeader()
	if err != nil {
		return nil, err
	}
	//payloadSignatureHeader := &cb.SignatureHeader{
	//	Creator: identity,
	//	Nonce:   nonce,
	//}

	paylBytes := protoutil.MarshalOrPanic(&cb.Payload{
		Header: protoutil.MakePayloadHeader(payloadChannelHeader, payloadSignatureHeader),
		Data:   data,
	})
	signature, err :=cry.Sign(paylBytes)
	//signature, err := c.Context().SigningManager().Sign(paylBytes, c.Context().PrivateKey())
	if err != nil {
		return nil, err
	}

	return &cb.Envelope{Payload: paylBytes, Signature: signature}, nil
}

// Receive receives events from the deliver server
func (c *DeliverConnection) Receive(eventch chan<- interface{}) {
	for {
		stream := c.deliverStream()
		if stream == nil {
			logger.Warn("The stream has closed. Terminating loop.")
			break
		}

		in, err := stream.Recv()

		logger.Debugf("Got deliver response: %#v", in)

		if c.Closed() {
			logger.Debugf("The connection has closed with error [%s]. Terminating loop.", err)
			break
		}

		if err == io.EOF {
			// This signifies that the stream has been terminated at the client-side. No need to send an event.
			logger.Debug("Received EOF from stream.")
			break
		}

		if err != nil {
			//logger.Warnf("Received error from stream: [%s]. Sending disconnected event.", err)
			eventch <- clientdsp.NewDisconnectedEvent(err)
			break
		}

		eventch <- NewEvent(in, c.url)
	}
	logger.Debug("Exiting stream listener")
}
