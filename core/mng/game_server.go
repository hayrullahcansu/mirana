package mng

import (
	"sync"

	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type GameServer struct {
	RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

// func CreateNewGameServer() *GameServer {
// roomManager := CreateNewRoomManager(50)
// lobby := CreateNewLobby()
// return &GameServer{
// 	RoomManager: roomManager,
// 	// Lobby:       lobby,
// 	// Users:       make(map[int]*dto.PlayerDto),
// }
// }

var _instance *GameServer

var _once sync.Once

func Instance() *GameServer {
	_once.Do(initialLogger)
	return _instance
}

func initialLogger() {
	// _instance = &config
}
