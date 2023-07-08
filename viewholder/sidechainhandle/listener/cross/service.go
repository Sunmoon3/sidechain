package cross

import (
	"fmt"
	"io"
	"viewholder/sidechainhandle/events"
	pb "viewholder/sidechainhandle/listener/cross/proto"
)

type CrossService struct {
	evch chan<- interface{}
}

func NewCrossService(evch chan<- interface{}) *CrossService {
	return &CrossService{
		evch: evch,
	}
}

func (cs *CrossService) CrossIn(stream pb.CrossService_CrossInServer) error {
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error while recv message: %v", err)
		}
		log.Debugf("recive crosslanguage %v", r)

		evt := events.NewCrossLanguageRecvedEvent(r)
		cs.evch <- evt

		resp := pb.CrossLanguageRespose{
			Status: 200,
			Msg:    "success recv the crosslanguage",
		}
		err = stream.Send(&resp)
		if err != nil {
			return fmt.Errorf("error while send message respose: %v", err)
		}
	}
}
