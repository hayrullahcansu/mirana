package sp

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/hayrullahcansu/mirana/core/comm/netsp"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
	"github.com/hayrullahcansu/mirana/core/mdl"
	"github.com/hayrullahcansu/mirana/core/types/gs"
	"github.com/hayrullahcansu/mirana/utils/que"
)

type SPGameRoom struct {
	*netw.BaseRoomManager
	Players     map[*netsp.NetSPClient]bool
	GameState   *gs.GameState
	GamePlayers []*netsp.SPPlayer
	// Pack        *list.List
	Pack *que.Queue
}

func NewSPGameRoom() *SPGameRoom {
	gameRoom := &SPGameRoom{
		Players:         make(map[*netsp.NetSPClient]bool),
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
		GamePlayers:     make([]*netsp.SPPlayer, 0, 6),
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
			fmt.Println("broadcast called")
			for client := range s.Players {
				select {
				case client.Send <- e:
				}
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
		default:

		}
	}
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
		client.AddMoney(addMoney.InternalId, addMoney.Amount)
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

func (m *SPGameRoom) OnDeal(c interface{}, deal *netw.Deal) {
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		client.Deal()
		m.Broadcast <- &netw.Envelope{
			Client: "client_id",
			Message: &netw.Deal{
				InternalId: deal.InternalId,
			},
			MessageCode: netw.EDeal,
		}
	}
	everyoneDealed := true
	for p, _ := range m.Players {
		if !p.IsDeal {
			everyoneDealed = false
		}
	}
	if everyoneDealed {
		go m.startGame()
	}
}

func (m *SPGameRoom) startGame() {
	m.Broadcast <- &netw.Envelope{
		Client: "server",
		Message: &netw.PlayGame{
			Mode: "game_will_start_in_3",
		},
		MessageCode: netw.EPlayGame,
	}
	initializeDone := make(chan bool, 1)

	go func() {
		m.init()
		var indexer int = 0

		for p1, _ := range m.Players {
			if len(p1.Players) > 0 && p1.IsDeal {
				for _, p := range p1.Players {
					m.GamePlayers = append(m.GamePlayers, p)
					indexer++
				}
			}
		}
		fmt.Printf("size : %d", len(m.GamePlayers))

		sort.Slice(m.GamePlayers, func(p, q int) bool {
			pp := m.GamePlayers[p]
			qq := m.GamePlayers[q]
			if pp == nil || qq == nil {
				fmt.Println("NİL GELDİ")
				return false
			}
			return m.GamePlayers[p].InternalId < m.GamePlayers[q].InternalId
		})
		fmt.Printf("size after sort : %d", len(m.GamePlayers))
		for _, val := range m.GamePlayers {
			card := m.PopCard()
			val.HitCard(card)
			fmt.Printf("   card:%s", card.String())
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.Hit{
					InternalId: val.InternalId,
					Card:       card.String(),
				},
				MessageCode: netw.EHit,
			}
			time.Sleep(time.Millisecond * 300)
		}
		initializeDone <- true
	}()
	time.Sleep(time.Second * 3)

	<-initializeDone
	close(initializeDone)

}

func (m *SPGameRoom) PopCard() *mdl.Card {
	element := m.Pack.Dequeue()
	if element != nil {
		return element.(*mdl.Card)
	}
	return nil
}

func (m *SPGameRoom) init() {
	m.Pack = que.Init()
	// var a = make([]interface{}, len(mdl.CardValues)*len(mdl.CardTypes)) // or slice := make([]int, elems)

	var a []*mdl.Card
	// var indexer = 0
	for _, cardValue := range mdl.CardValues {
		for _, cardType := range mdl.CardTypes {
			c := mdl.NewCardData(cardType, cardValue)
			a = append(a, c)
			// a[indexer] =
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	for _, v := range a {
		m.Pack.Enqueue(v)
	}
}
