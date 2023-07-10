/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "git.labpc.bluarry.top/bluarry/viewholder/chain"

// Connection defines the functions for an event server connection
type Connection interface {
	// Receive sends events to the given channel
	Receive(chan<- interface{})
	// Close closes the connection
	Close()
	// Closed return true if the connection is closed
	Closed() bool
}

// ConnectionProvider creates a Connection.
//type ConnectionProvider func(context context.Client, chConfig fab.ChannelCfg, peer fab.Peer) (Connection, error)
type ConnectionProvider func( *chain.Config, *chain.Node) (Connection, error)