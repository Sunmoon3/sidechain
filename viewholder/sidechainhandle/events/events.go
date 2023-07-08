package events

import (
	pb "viewholder/sidechainhandle/listener/cross/proto"
)

type Event interface{}

type CrossTransferEvent struct {
	Txid    string
	Payload string
}

func NewCrossTransferEvent(txid, payload string) *CrossTransferEvent {
	return &CrossTransferEvent{
		Txid:    txid,
		Payload: payload,
	}
}

type StopEvent struct {
	ErrCh chan<- error
}

func NewStopEvent(errch chan<- error) *StopEvent {
	return &StopEvent{
		ErrCh: errch,
	}
}

func (cte *CrossTransferEvent) GettxId() string {
	return cte.Txid
}

func (cte *CrossTransferEvent) GetPayload() string {
	return cte.Payload
}

type CrosslanguageGeneratedEvent struct {
	Req *pb.CrossLanguageRequest
}

func NewCrosslanguageGeneratedEvent(req *pb.CrossLanguageRequest) *CrosslanguageGeneratedEvent {
	return &CrosslanguageGeneratedEvent{
		Req: req,
	}
}

type CrossLanguageRecvedEvent struct {
	Req *pb.CrossLanguageRequest
}

func NewCrossLanguageRecvedEvent(req *pb.CrossLanguageRequest) *CrossLanguageRecvedEvent {
	return &CrossLanguageRecvedEvent{
		Req: req,
	}
}
