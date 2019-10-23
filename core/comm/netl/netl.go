package netl

import (
	"fmt"

	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type NetLobbyClient struct {
	netw.BaseClient
	netw.EnvelopeListener
}

func (c *NetLobbyClient) OnNotify(notify *netw.Notify) {
	fmt.Println("WORKED INHERITED METHOD")
}

func NewClient() *NetLobbyClient {
	return &NetLobbyClient{
		BaseClient: netw.BaseClient{
			Send:       make(chan *netw.Envelope),
			Notify:     make(chan *netw.Notify),
			Unregister: make(chan interface{}),
		},
	}
}
