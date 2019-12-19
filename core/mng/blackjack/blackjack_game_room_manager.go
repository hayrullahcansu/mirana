package blackjack

import (
	"sync"
	"time"

	"bitbucket.org/digitdreamteam/mirana/core/settings"
)

type BlackjackGameRoomManager struct {
	GameRooms   map[*BlackjackGameRoom]bool
	DefaultRoom *BlackjackGameRoom
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
}

var _instance *BlackjackGameRoomManager

var _once sync.Once

func Manager() *BlackjackGameRoomManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &BlackjackGameRoomManager{
		GameRooms:   make(map[*BlackjackGameRoom]bool),
		DefaultRoom: NewBlackjackGameRoom(SINGLE_DECK),
	}
	go func() {
		_instance.work()
	}()
}

func (manager *BlackjackGameRoomManager) work() {
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

func (manager *BlackjackGameRoomManager) requestPlayGame(c *BlackjackClient, gameType GameType) {
	g := NewBlackjackGameRoom(gameType)
	manager.GameRooms[g] = true
	g.ConnectGame(c)
}

func (manager *BlackjackGameRoomManager) RequestSingleDeckPlayGame(c *BlackjackClient) {
	manager.requestPlayGame(c, SINGLE_DECK)
}

func (manager *BlackjackGameRoomManager) RequestAmericanPlayGame(c *BlackjackClient) {
	manager.requestPlayGame(c, AMERICAN)
}

func (manager *BlackjackGameRoomManager) RemoveGameRoom(r *BlackjackGameRoom) {
	_, ok := manager.GameRooms[r]
	if ok {
		delete(manager.GameRooms, r)
	}
}
