package exroom

import (
	"fmt"
	"sync"

	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
)

type ExroomManager struct {
	*netw.BaseRoomManager
	Players map[*netw.BaseClient]bool
}

var _instance *ExroomManager

var _once sync.Once

func Manager() *ExroomManager {
	_once.Do(initialManagerInstance)
	return _instance
}

func initialManagerInstance() {
	_instance = &ExroomManager{
		BaseRoomManager: netw.NewBaseRoomManager(),
		Players:         make(map[*netw.BaseClient]bool),
	}
	go _instance.ListenEvents()
}

func (s *ExroomManager) ListenEvents() {
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
			// default:
			// 	break
		}
	}
}

func (s *ExroomManager) OnNotify(notify *netw.Notify) {
	switch v := notify.Message.Message.(type) {
	case netw.Event:
	case netw.Stamp:
	case netw.AddMoney:
	case netw.Deal:
	case netw.Stand:
	case netw.Hit:
	case netw.Double:
	case netw.PlayGame:
		s.OnPlayGame(notify.SentBy, notify.Message.Message.(*netw.PlayGame))
	default:
		fmt.Printf("unexpected type %T", v)
	}
}
func (m *ExroomManager) ConnectLobby(client *netw.BaseClient) {
	m.Players[client] = true
	client.Notify = m.Notify
	m.Register <- client
}

func (m *ExroomManager) OnConnect(client interface{}) {
	_, ok := client.(*netw.BaseClient)
	if ok {

	}
}
