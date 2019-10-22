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
	// Users       map[int]*dto.PlayerDto
}

var _instance *GameServer

var _once sync.Once

func Instance() *GameServer {
	_once.Do(initialLogger)
	return _instance
}

func initialLogger() {
	// _instance = &config
}
