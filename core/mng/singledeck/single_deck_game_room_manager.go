package singledeck

import (
	"sync"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/settings"
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
	go func() {
		_instance.work()
	}()
}

func (manager *SingleDeckGameRoomManager) work() {
	purgeTicker := time.NewTicker(settings.PURGE_PERIOD)
	echoTicker := time.NewTicker(settings.ECHO_PERIOD)
	defer func() {
		purgeTicker.Stop()
		echoTicker.Stop()
	}()
	for {
		select {
		case <-purgeTicker.C:
			for _, ok := range _instance.GameRooms {
				if ok {

				}
			}
		case <-echoTicker.C:
			for r, ok := range _instance.GameRooms {
				if ok {
					r.PrintRoomStatus()
				}
			}
		}
	}
}

func (manager *SingleDeckGameRoomManager) RequestPlayGame(c *SingleDeckSPClient) {
	g := NewSingleDeckGameRoom()
	manager.GameRooms[g] = true
	g.ConnectGame(c)
}

func (manager *SingleDeckGameRoomManager) RemoveGameRoom(r *SingleDeckGameRoom) {
	_, ok := manager.GameRooms[r]
	if ok {
		delete(manager.GameRooms, r)
	}
}
