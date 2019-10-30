package lobby2

import (
	"fmt"
	"sync"

	"github.com/hayrullahcansu/mirana/core/types"
)

type LobbyManager struct {
	types.EnvelopeListener
	Players    map[*types.Player]bool
	Register   chan *types.Player
	Unregister chan *types.Player
	Notify     chan *types.Notify
	Broadcast  chan *types.Envelope
}

var _instance *LobbyManager

var _once sync.Once

func Manager() *LobbyManager {
	_once.Do(initialGameManagerInstance)
	return _instance
}

func initialGameManagerInstance() {
	_instance = &LobbyManager{
		Players:    make(map[*types.Player]bool),
		Register:   make(chan *types.Player),
		Unregister: make(chan *types.Player),
		Notify:     make(chan *types.Notify),
		Broadcast:  make(chan *types.Envelope),
	}

}

func (manager *LobbyManager) RequestPlayGame(client *types.Client) {
	player := &types.Player{
		Client:     client,
		InternalID: -1,
		Unregister: manager.Unregister,
	}
	manager.Players[player] = true
	client.Notify = manager.Notify
	manager.Register <- player
}

func (s *LobbyManager) listenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			if _, ok := s.Players[player]; ok {
				// s.OnLeave(player)
			}
		case message := <-s.Broadcast:
			for player := range s.Players {
				select {
				case player.Client.Send <- message:
					// default:
					// 	s.Unregister <- player
				}
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
			// default:
			// 	break
		}
	}
}

func (s *LobbyManager) OnConnect(player *types.Player) {}

func (s *LobbyManager) OnNotify(notify *types.Notify) {
	switch v := notify.Message.Message.(type) {
	case types.Event:
	case types.Stamp:
	case types.AddMoney:
	case types.Deal:
	case types.Stand:
	case types.Hit:
	case types.Double:
	default:
		fmt.Printf("unexpected type %T", v)
	}
}
func (s *LobbyManager) OnEvent(event *types.Event)          {}
func (s *LobbyManager) OnStamp(stamp *types.Stamp)          {}
func (s *LobbyManager) OnAddMoney(addMoney *types.AddMoney) {}
func (s *LobbyManager) OnDeal(deal *types.Deal)             {}
func (s *LobbyManager) OnStand(stand *types.Stand)          {}
func (s *LobbyManager) OnHit(hit *types.Hit)                {}
func (s *LobbyManager) OnDouble(double *types.Double)       {}
func (s *LobbyManager) OnPlayGame(playGame *types.PlayGame) {}
