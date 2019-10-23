package netsp

import (
	"fmt"

	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type NetSPClient struct {
	netw.BaseClient
	netw.EnvelopeListener
}

func (c *NetSPClient) OnNotify(notify *netw.Notify) {
	fmt.Println("WORKED INHERITED METHOD")
}

func NewClient() *NetSPClient {
	return &NetSPClient{
		BaseClient: netw.BaseClient{},
	}
}
