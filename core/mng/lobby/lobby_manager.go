package lobby

import (
	"sync"

	"github.com/hayrullahcansu/mirana/core/types"
)

type LobbyManager struct {
	// GameRooms map[*SPGameRoom]bool
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *LobbyManager

var _once sync.Once

func Manager() *LobbyManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &LobbyManager{
		// GameRooms: make(map[*SPGameRoom]bool),
	}
}

func (manager *LobbyManager) RequestPlayGame(client *types.Client) {
	player := &types.Player{
		Client:     client,
		InternalID: -1,
		Unregister: gameRoom.Unregister,
	}
	manager.GameRooms[gameRoom] = true
	client.notify = gameRoom.Notify
	gameRoom.register <- player
}
