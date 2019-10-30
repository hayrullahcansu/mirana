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
		Register:   make(chan interface{}, 1),
		Unregister: make(chan interface{}, 1),
		Notify:     make(chan *Notify, 1),
		Broadcast:  make(chan *Envelope, 1),
	}
}

func (s *BaseRoomManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case b := <-s.Broadcast:
			println(b)
			break
		case notify := <-s.Notify:
			s.OnNotify(notify)
			// default:
			// break
		}
	}
}

func (m *BaseRoomManager) OnConnect(c interface{}) {

}

func (m *BaseRoomManager) OnDisconnect(c interface{}) {

}

func (s *BaseRoomManager) OnNotify(notify *Notify) {
	d := notify.Message.Message
	switch v := notify.Message.Message.(type) {
	case Event:
		t := d.(Event)
		s.OnEvent(notify.SentBy, &t)
	case Stamp:
		t := d.(Stamp)
		s.OnStamp(notify.SentBy, &t)
	case AddMoney:
		t := d.(AddMoney)
		s.OnAddMoney(notify.SentBy, &t)
	case Deal:
		t := d.(Deal)
		s.OnDeal(notify.SentBy, &t)
	case Stand:
		t := d.(Stand)
		s.OnStand(notify.SentBy, &t)
	case Hit:
		t := d.(Hit)
		s.OnHit(notify.SentBy, &t)
	case Double:
		t := d.(Double)
		s.OnDouble(notify.SentBy, &t)
	case PlayGame:
		t := d.(PlayGame)
		s.OnPlayGame(notify.SentBy, &t)
	default:
		fmt.Printf("unexpected type %T", v)
	}
}

func (s *BaseRoomManager) OnEvent(c interface{}, event *Event)          {}
func (s *BaseRoomManager) OnStamp(c interface{}, stamp *Stamp)          {}
func (s *BaseRoomManager) OnAddMoney(c interface{}, addMoney *AddMoney) {}
func (s *BaseRoomManager) OnDeal(c interface{}, deal *Deal)             {}
func (s *BaseRoomManager) OnStand(c interface{}, stand *Stand)          {}
func (s *BaseRoomManager) OnHit(c interface{}, hit *Hit)                {}
func (s *BaseRoomManager) OnDouble(c interface{}, double *Double)       {}
func (s *BaseRoomManager) OnPlayGame(c interface{}, playGame *PlayGame) {
	fmt.Printf("expo type %T", playGame)
}
