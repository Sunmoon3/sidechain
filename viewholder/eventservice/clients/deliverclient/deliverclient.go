package deliverclient

import (
	"git.labpc.bluarry.top/bluarry/viewholder/chain"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/api"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/client"
	deliverconn "git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/deliverclient/connection"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/deliverclient/dispatcher"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/clients/deliverclient/seek"
	"git.labpc.bluarry.top/bluarry/viewholder/eventservice/common"
	ab "github.com/hyperledger/fabric-protos-go/orderer"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math"
	"time"
)

var deliverProvider=func(config *chain.Config,node *chain.Node)(api.Connection,error){
	if config==nil||node==nil{
		return nil,errors.New("config is nil")
	}

	return deliverconn.New(config,node,deliverconn.Deliver)
}


var deliverFilteredProvider=func(config *chain.Config,node *chain.Node)(api.Connection,error){
	if config==nil||node==nil{
		return nil,errors.New("config is nil")
	}

	return deliverconn.New(config,node,deliverconn.DeliverFiltered)
}


var logger=log.New()
type Client struct {
	params
	*client.Client
}


func New(confpath string,opts ...common.Opt)(*Client ,error){

	params := defaultParams()
	common.Apply(params,opts)

	cfg,err:=chain.LoadConfig(confpath)
	if err != nil {
		return nil, err
	}

	dispa:=dispatcher.New(&cfg, params.connProvider,opts...)



	if params.seekType == "" {
		params.seekType = seek.Newest

		//discard (do not publish) next BlockEvent/FilteredBlockEvent in dispatcher, since default seek type 'newest' is
		// only needed for block height calculations
		dispa.UpdateLastBlockInfoOnly()
	}

	cli:=&Client{
		params: *params,
		Client: client.New(dispa),
	}

	cli.SetAfterConnectHandler(cli.seek)
	cli.SetBeforeReconnectHandler(cli.setSeekFromLastBlockReceived)

	if err := cli.Start(); err != nil {
		return nil, err
	}

	return cli,nil
}



func (c *Client) seek() error {
	logger.Debug("Sending seek request....")

	seekInfo, err := c.seekInfo()
	if err != nil {
		return err
	}

	errch := make(chan error, 1)
	err1 := c.Submit(dispatcher.NewSeekEvent(seekInfo, errch))
	if err1 != nil {
		return err1
	}
	select {
	case err = <-errch:
	case <-time.After(c.respTimeout):
		err = errors.New("timeout waiting for deliver status response")
	}

	if err != nil {
		logger.Errorf("Unable to send seek request: %s", err)
		return err
	}

	logger.Debug("Successfully sent seek")
	return nil
}



func (c *Client) setSeekFromLastBlockReceived() error {
	c.Lock()
	defer c.Unlock()

	// Make sure that, when we reconnect, we receive all of the events that we've missed
	lastBlockNum := c.Dispatcher().LastBlockNum()
	if lastBlockNum < math.MaxUint64 {
		c.seekType = seek.FromBlock
		c.fromBlock = c.Dispatcher().LastBlockNum() + 1
		logger.Debugf("Setting seek info from last block received + 1: %d", c.fromBlock)
	} else {
		// We haven't received any blocks yet. Just ask for the newest
		logger.Debugf("Setting seek info from newest")
		c.seekType = seek.Newest
	}
	return nil
}




func (c *Client) seekInfo() (*ab.SeekInfo, error) {
	c.RLock()
	defer c.RUnlock()

	switch c.seekType {
	case seek.Newest:
		logger.Debugf("Returning seek info: Newest")
		return seek.InfoNewest(), nil
	case seek.Oldest:
		logger.Debugf("Returning seek info: Oldest")
		return seek.InfoOldest(), nil
	case seek.FromBlock:
		logger.Debugf("Returning seek info: FromBlock(%d)", c.fromBlock)
		return seek.InfoFrom(c.fromBlock), nil
	default:
		return nil, errors.Errorf("unsupported seek type:[%s]", c.seekType)
	}
}
