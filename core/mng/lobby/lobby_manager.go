package lobby

import (
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/hayrullahcansu/mirana/core/comm/netw"
)

type LobbyManager struct {
	*netw.BaseRoomManager
	Players map[*NetLobbyClient]bool
}

var _instance *LobbyManager

var _once sync.Once

func Manager() *LobbyManager {
	_once.Do(initialManagerInstance)
	return _instance
}

func initialManagerInstance() {
	_instance = &LobbyManager{
		BaseRoomManager: netw.NewBaseRoomManager(),
		Players:         make(map[*NetLobbyClient]bool),
	}
	go _instance.ListenEvents()
}

func (s *LobbyManager) ListenEvents() {
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case _ = <-s.Broadcast:
		case notify := <-s.Notify:
			s.OnNotify(notify)
			// default:
		}
	}
}

func (s *LobbyManager) OnNotify(notify *netw.Notify) {
	d := notify.Message.Message
	switch v := notify.Message.Message.(type) {
	case netw.Event:
		t := d.(netw.Event)
		s.OnEvent(notify.SentBy, &t)
	case netw.Stamp:
		t := d.(netw.Stamp)
		s.OnStamp(notify.SentBy, &t)
	case netw.Split:
		t := d.(netw.Split)
		s.OnSplit(notify.SentBy, &t)
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

func (m *LobbyManager) ConnectLobby(c *NetLobbyClient) {
	m.Players[c] = true
	c.Notify = m.Notify
	m.Register <- c
}

func (m *LobbyManager) OnConnect(c interface{}) {
	client, ok := c.(*NetLobbyClient)
	if ok {
		client.Unregister = m.Unregister
	}
}
func (m *LobbyManager) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	client, ok := c.(*NetLobbyClient)
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
