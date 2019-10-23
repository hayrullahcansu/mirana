package lobby

import (
	"sync"

	"github.com/hayrullahcansu/mirana/core/comm/netl"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type LobbyManager struct {
	*netw.BaseRoomManager
	Players map[*netl.NetLobbyClient]bool
}

var _instance *LobbyManager

var _once sync.Once

func Manager() *LobbyManager {
	_once.Do(initialManagerInstance)
	return _instance
}

func initialManagerInstance() {
	_instance = &LobbyManager{
		BaseRoomManager: netw.NewBaseRoomManager(),
		Players:         make(map[*netl.NetLobbyClient]bool),
	}
	go _instance.ListenEvents()
}

func (s *LobbyManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case _ = <-s.Broadcast:
			break
		case notify := <-s.Notify:
			s.OnNotify(notify)
		default:
			break
		}
	}
}
func (m *LobbyManager) ConnectLobby(client *netl.NetLobbyClient) {
	m.Players[client] = true
	client.Notify = m.Notify
	m.Register <- client
}

func (m *LobbyManager) OnConnect(client interface{}) {
	_, ok := client.(*netl.NetLobbyClient)
	if ok {

	}
}
