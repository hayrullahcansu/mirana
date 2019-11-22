package singledeck

import (
	"sync"
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

func (manager *SingleDeckGameRoomManager) RequestPlayGame(c *SingleDeckSPClient) {
	g := NewSingleDeckGameRoom()
	manager.GameRooms[g] = true
	g.ConnectGame(c)
}
