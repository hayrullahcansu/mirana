package exroom

import (
	"sync"

	"github.com/hayrullahcansu/mirana/core/comm/netl"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type ExroomManager struct {
	*netw.BaseRoomManager
	Players map[*netl.NetLobbyClient]bool
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
		Players:         make(map[*netl.NetLobbyClient]bool),
	}
	go _instance.ListenEvents()
}
