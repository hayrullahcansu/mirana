package netw

import (
	"fmt"
)

type BaseRoomManager struct {
	EnvelopeListener
	Register   chan interface{}
	Unregister chan interface{}
	Notify     chan *Notify
	Broadcast  chan *Envelope
}

type IBaseRoomManager interface {
	OnConnect(baseClient *BaseClient)
	OnDisconnect(baseClient *BaseClient)
}

func NewBaseRoomManager() *BaseRoomManager {
	return &BaseRoomManager{
		Register:   make(chan interface{}),
		Unregister: make(chan interface{}),
		Notify:     make(chan *Notify),
		Broadcast:  make(chan *Envelope),
	}
}

func (s *BaseRoomManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case _ = <-s.Broadcast:
			break
		case notify := <-s.Notify:
			s.OnNotify(notify)
		default:
			break
		}
	}
}

func (m *BaseRoomManager) OnConnect(client interface{}) {

}

func (m *BaseRoomManager) OnDisconnect(client interface{}) {

}

func (s *BaseRoomManager) OnNotify(notify *Notify) {
	switch v := notify.Message.Message.(type) {
	case Event:
	case Stamp:
	case AddMoney:
	case Deal:
	case Stand:
	case Hit:
	case Double:
	default:
		fmt.Printf("unexpected type %T", v)
	}
}

func (s *BaseRoomManager) OnEvent(event *Event)          {}
func (s *BaseRoomManager) OnStamp(stamp *Stamp)          {}
func (s *BaseRoomManager) OnAddMoney(addMoney *AddMoney) {}
func (s *BaseRoomManager) OnDeal(deal *Deal)             {}
func (s *BaseRoomManager) OnStand(stand *Stand)          {}
func (s *BaseRoomManager) OnHit(hit *Hit)                {}
func (s *BaseRoomManager) OnDouble(double *Double)       {}
func (s *BaseRoomManager) OnPlayGame(playGame *PlayGame) {}
