package netsp

import (
	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type NetSPClient struct {
	*netw.BaseClient
	amount float32
}

func NewClient() *NetSPClient {

	client := &NetSPClient{}
	base := netw.NewBaseClient(client)
	client.BaseClient = base
	return client
}

func (c *NetSPClient) AddMoney(amount float32) {
	c.amount += amount
}
