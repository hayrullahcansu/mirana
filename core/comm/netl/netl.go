package netl

import (
	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type NetLobbyClient struct {
	*netw.BaseClient
}

func NewClient() *NetLobbyClient {
	return &NetLobbyClient{
		BaseClient: &netw.BaseClient{
			Send:       make(chan *netw.Envelope),
			Notify:     make(chan *netw.Notify),
			Unregister: make(chan interface{}),
		},
	}
}
