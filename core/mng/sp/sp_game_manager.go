package sp

import "sync"

type SPGameManager struct {
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *SPGameManager

var _once sync.Once

func Manager() *SPGameManager {
	_once.Do(initialLogger)
	return _instance
}

func initialLogger() {
	_instance = &SPGameManager{}
}
func (manager *SPGameManager) RequestPlayGame(client *Client) {

}
