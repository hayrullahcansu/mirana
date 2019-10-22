package sp

import (
	"container/list"
	"fmt"

	"github.com/hayrullahcansu/mirana/core/types"
	"github.com/hayrullahcansu/mirana/core/types/gs"
)

type SPGameRoom struct {
	Pack      *list.List
	Players   map[*types.Player]bool
	GameState *gs.GameState
	// RoomManager *RoomManager
	// Lobby       *Lobby
	// Users       map[int]*dto.PlayerDto
	Broadcast  chan *types.Envelope
	Notify     chan *types.Notify
	unregister chan *types.Player
	register   chan *types.Player
}

func NewSPGameRoom() *SPGameRoom {
	gameRoom := &SPGameRoom{
		Players:   make(map[*types.Player]bool),
		GameState: gs.NewGameState(),
		Notify:    make(chan *types.Notify),
	}
	return gameRoom
}

func (s *SPGameRoom) ListenEvents() {
	for {
		select {
		case player := <-s.register:
			s.OnConnect(player)
		case player := <-s.unregister:
			if _, ok := s.Players[player]; ok {
				s.OnLeave(player)
			}
		case message := <-s.Broadcast:
			for player := range s.Players {
				select {
				case player.Client.Send <- message:
				default:
					s.unregister <- player
				}
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
		default:
		}
	}
}
func (s *SPGameRoom) OnConnect(player *types.Player) {
	s.Players[player] = true
}

func (s *SPGameRoom) OnLeave(player *types.Player) {
	delete(s.Players, player)
	close(player.Client.Send)
}
func (s *SPGameRoom) OnNotify(notify *types.Notify) {
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
func (s *SPGameRoom) PopCard() *types.Card {
	element := s.Pack.Front().Value
	if element != nil {
		return element.(*types.Card)
	}
	return nil
}

func (s *SPGameRoom) NextGameStatu() {
	_, _new := s.GameState.GetNextStatu()
	switch _new {
	case gs.INIT:
	case gs.WAIT_PLAYERS:
	case gs.DEALING:
	case gs.PREPARING:
	case gs.IN_PLAY:
	case gs.DONE:
	default:
	}
}

func (s *SPGameRoom) init() {
	s.Pack = list.New()
	for _, cardValue := range types.CardValues {
		for _, cardType := range types.CardTypes {
			card := types.NewCardData(cardType, cardValue)
			s.Pack.PushBack(card)
		}
	}
}
