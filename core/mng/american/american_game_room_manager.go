package american

import (
	"sync"
)

type AmericanGameRoomManager struct {
	GameRooms   map[*AmericanGameRoom]bool
	DefaultRoom *AmericanGameRoom
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *AmericanGameRoomManager

var _once sync.Once

func Manager() *AmericanGameRoomManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &AmericanGameRoomManager{
		GameRooms:   make(map[*AmericanGameRoom]bool),
		DefaultRoom: NewAmericanGameRoom(),
	}
}

func (manager *AmericanGameRoomManager) RequestPlayGame(c *AmericanSPClient) {
	g := NewAmericanGameRoom()
	manager.GameRooms[g] = true
	g.ConnectGame(c)
}
