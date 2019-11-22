package lobby

import (
	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
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
