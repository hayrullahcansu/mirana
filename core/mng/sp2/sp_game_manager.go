package sp2

// import (
// 	"sync"

// 	"bitbucket.org/digitdreamteam/mirana/core/types"
// )

// type SPGameManager struct {
// 	GameRooms map[*SPGameRoom]bool
// 	// RoomManager *RoomManager
// 	// Lobby       *Lobby
// 	// Users       map[int]*dto.PlayerDto
// }

// var _instance *SPGameManager

// var _once sync.Once

// func Manager() *SPGameManager {
// 	_once.Do(initialGameManagerInstance)
// 	return _instance
// }

// func initialGameManagerInstance() {
// 	_instance = &SPGameManager{
// 		GameRooms: make(map[*SPGameRoom]bool),
// 	}
// }

// func (manager *SPGameManager) RequestPlayGame(client *types.Client) {
// 	gameRoom := NewSPGameRoom()
// 	player := &types.Player{
// 		Client:     client,
// 		InternalID: -1,
// 		Unregister: gameRoom.unregister,
// 	}
// 	manager.GameRooms[gameRoom] = true
// 	client.notify = gameRoom.Notify
// 	gameRoom.register <- player
// }
