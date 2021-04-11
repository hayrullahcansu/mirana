package lobby

import (
	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type NetLobbyClient struct {
	*netw.BaseClient
}

func NewClient() *NetLobbyClient {
	client := &NetLobbyClient{}
	base := netw.NewBaseClient(client)
	client.BaseClient = base
	return client
}
