package cross

import (
	"fmt"
	"viewholder/config"

	"google.golang.org/grpc"
	"net"
	pb "viewholder/sidechainhandle/listener/cross/proto"
)

var log = logger.GetLogger(config.CROSS_LISTENER)

type CrossListener struct {
	errch  chan<- error
	evch   chan<- interface{}
	lis    *net.Listener
	server *grpc.Server
}

func NewCin(errch chan<- error) *CrossListener {
	cl := &CrossListener{
		errch: errch,
	}
	cl.receiveFromMC(errch)
	return cl
}

func (cl *CrossListener) Start(evch chan<- interface{}) {
	log.Info("CrossListener Start...")
	cl.evch = evch
	crossService := NewCrossService(evch)
	pb.RegisterCrossServiceServer(cl.server, crossService)

	go func() {
		log.Infof("The CorssChain Listener has been Started..")
		log.Infof("The CorssChain Listener listening at %v ....", (*cl.lis).Addr())
		err := cl.server.Serve(*cl.lis)
		if err != nil {
			cl.errch <- fmt.Errorf("error while start listen: %v", err)
			return
		}
	}()
}

func (cl *CrossListener) Stop() {
	cl.server.Stop()
}
