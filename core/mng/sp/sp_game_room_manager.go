package sp

import (
	"sync"

	"github.com/hayrullahcansu/mirana/core/comm/netsp"
)

type SPGameRoomManager struct {
	GameRooms   map[*SPGameRoom]bool
	DefaultRoom *SPGameRoom
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *SPGameRoomManager

var _once sync.Once

func Manager() *SPGameRoomManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &SPGameRoomManager{
		GameRooms:   make(map[*SPGameRoom]bool),
		DefaultRoom: NewSPGameRoom(),
	}
}

func (manager *SPGameRoomManager) RequestPlayGame(c *netsp.NetSPClient) {
	g := NewSPGameRoom()
	manager.GameRooms[g] = true
	g.ConnectGame(c)
	// manager.DefaultRoom.ConnectGame(c)

	// m.Players[c] = true
	// c.Notify = m.Notify
	// m.Register <- c
}