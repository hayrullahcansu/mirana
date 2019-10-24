package sp

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/hayrullahcansu/mirana/core/comm/netsp"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
	"github.com/hayrullahcansu/mirana/core/types/gs"
)

type SPGameRoom struct {
	*netw.BaseRoomManager
	Players   map[*netsp.NetSPClient]bool
	GameState *gs.GameState
}

func NewSPGameRoom() *SPGameRoom {
	gameRoom := &SPGameRoom{
		Players:         make(map[*netsp.NetSPClient]bool),
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
	}
	go gameRoom.ListenEvents()
	return gameRoom
}

func (s *SPGameRoom) ListenEvents() {
	fmt.Println("GIRDI")
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case e := <-s.Broadcast:
			for client := range s.Players {
				select {
				case client.Send <- e:
				default:
					client.Unregister <- client
				}
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
		default:

		}
	}
	fmt.Println("CIKTI")
}

func (s *SPGameRoom) OnNotify(notify *netw.Notify) {
	d := notify.Message.Message
	switch v := notify.Message.Message.(type) {
	case netw.Event:
		t := d.(netw.Event)
		s.OnEvent(notify.SentBy, &t)
	case netw.Stamp:
		t := d.(netw.Stamp)
		s.OnStamp(notify.SentBy, &t)
	case netw.AddMoney:
		t := d.(netw.AddMoney)
		s.OnAddMoney(notify.SentBy, &t)
	case netw.Deal:
		t := d.(netw.Deal)
		s.OnDeal(notify.SentBy, &t)
	case netw.Stand:
		t := d.(netw.Stand)
		s.OnStand(notify.SentBy, &t)
	case netw.Hit:
		t := d.(netw.Hit)
		s.OnHit(notify.SentBy, &t)
	case netw.Double:
		t := d.(netw.Double)
		s.OnDouble(notify.SentBy, &t)
	case netw.PlayGame:
		t := d.(netw.PlayGame)
		s.OnPlayGame(notify.SentBy, &t)
	default:
		fmt.Printf("unexpected type %T", v)
	}
}

func (m *SPGameRoom) ConnectGame(c *netsp.NetSPClient) {
	m.Players[c] = true
	c.Notify = m.Notify
	m.Register <- c
}

func (m *SPGameRoom) OnConnect(c interface{}) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {

	}
}
func (m *SPGameRoom) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		//TODO: check player able to play?
		mode := playGame.Mode
		guid := uuid.New()
		playGame.Id = guid.String()
		playGame.Mode = mode
		client.Send <- &netw.Envelope{
			Client:      "client_id",
			Message:     playGame,
			MessageCode: netw.EPlayGame,
		}
	}
}

func (m *SPGameRoom) OnAddMoney(c interface{}, addMoney *netw.AddMoney) {
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		client.AddMoney(addMoney.Amount)
		m.Broadcast <- &netw.Envelope{
			Client: "client_id",
			Message: &netw.AddMoney{
				InternalId: addMoney.InternalId,
				Amount:     addMoney.Amount,
			},
			MessageCode: netw.EAddMoney,
		}
	}
}
