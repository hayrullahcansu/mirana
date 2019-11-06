package singledeck

import (
	"sync"

	"bitbucket.org/digitdreamteam/mirana/core/comm/netsp"
)

type SingleDeckGameRoomManager struct {
	GameRooms   map[*SingleDeckGameRoom]bool
	DefaultRoom *SingleDeckGameRoom
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *SingleDeckGameRoomManager

var _once sync.Once

func Manager() *SingleDeckGameRoomManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &SingleDeckGameRoomManager{
		GameRooms:   make(map[*SingleDeckGameRoom]bool),
		DefaultRoom: NewSingleDeckGameRoom(),
	}
}

func (manager *SingleDeckGameRoomManager) RequestPlayGame(c *netsp.NetSPClient) {
	g := NewSingleDeckGameRoom()
	manager.GameRooms[g] = true
	g.ConnectGame(c)
	// manager.DefaultRoom.ConnectGame(c)

	// m.Players[c] = true
	// c.Notify = m.Notify
	// m.Register <- c
}
