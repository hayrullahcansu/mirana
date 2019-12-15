package american

import (
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/settings"
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
	go func() {
		_instance.work()
	}()
}

func (manager *AmericanGameRoomManager) work() {
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

func (manager *AmericanGameRoomManager) RequestPlayGame(c *AmericanSPClient) {
	g := NewAmericanGameRoom()
	logrus.Infof("American Game Room Created")
	manager.GameRooms[g] = true
	g.ConnectGame(c)
}

func (manager *AmericanGameRoomManager) RemoveGameRoom(r *AmericanGameRoom) {
	_, ok := manager.GameRooms[r]
	if ok {
		delete(manager.GameRooms, r)
	}
}
