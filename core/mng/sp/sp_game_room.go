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
	Players             map[*netsp.NetSPClient]bool
	GameState           *gs.GameState
	GamePlayers         []*netsp.SPPlayer
	System              *netsp.SPPlayer
	GameStateEvent      chan gs.GameStatu
	CurrentPlayerCursor int
	TurnOfPlay          string
	Pack                *que.Queue
}

func NewSPGameRoom() *SPGameRoom {
	gameRoom := &SPGameRoom{
		Players:         make(map[*netsp.NetSPClient]bool),
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
		GamePlayers:     make([]*netsp.SPPlayer, 0, 6),
		GameStateEvent:  make(chan gs.GameStatu, 1),
	}
	go gameRoom.ListenEvents()
	gameRoom.GameStateEvent <- gs.INIT
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
				}
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
		// default:
		case gameStateEvent := <-s.GameStateEvent:
			s.gameStateChanged(gameStateEvent)
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
		m.GameStateEvent <- gs.PREPARING
	}
}

func (m *SPGameRoom) startGame() {
	m.System = netsp.NewSPSystemPlayer()
	m.Broadcast <- &netw.Envelope{
		Client: "server",
		Message: &netw.PlayGame{
			Mode: "game_will_start_in_3",
		},
		MessageCode: netw.EPlayGame,
	}
	initializeDone := make(chan bool, 1)

	go func() {
		var indexer int = 0

		for p1, _ := range m.Players {
			if len(p1.Players) > 0 && p1.IsDeal {
				for _, p := range p1.Players {
					m.GamePlayers = append(m.GamePlayers, p)
					indexer++
				}
			}
		}

		sort.Slice(m.GamePlayers, func(p, q int) bool {
			pp := m.GamePlayers[p]
			qq := m.GamePlayers[q]
			if pp == nil || qq == nil {
				fmt.Println("NİL GELDİ")
				return false
			}
			return m.GamePlayers[p].InternalId < m.GamePlayers[q].InternalId
		})
		for _, player := range m.GamePlayers {
			m.pull_card_for_player(player)
			time.Sleep(time.Millisecond * 300)
		}
		m.pull_card_for_system()
		for _, player := range m.GamePlayers {
			m.pull_card_for_player(player)
			time.Sleep(time.Millisecond * 300)
		}
		m.pull_card_for_system()
		m.CurrentPlayerCursor = len(m.GamePlayers) - 1
		initializeDone <- true
	}()
	time.Sleep(time.Second * 3)

	<-initializeDone
	close(initializeDone)
	m.GameStateEvent <- gs.IN_PLAY

}

// func (m *SPGameRoom) OnEvent(c interface{}, event *netw.Event) {
// 	client, ok := c.(*netsp.NetSPClient)
// 	if ok {
// 		// client.AddMoney(addMoney.InternalId, addMoney.Amount)
// 		// m.Broadcast <- &netw.Envelope{
// 		// 	Client: "client_id",
// 		// 	Message: &netw.AddMoney{
// 		// 		InternalId: addMoney.InternalId,
// 		// 		Amount:     addMoney.Amount,
// 		// 	},
// 		// 	MessageCode: netw.EAddMoney,
// 		// }
// 	}
// }

func (m *SPGameRoom) OnHit(c interface{}, hit *netw.Hit) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {
		for _, player := range m.GamePlayers {
			if player.InternalId == hit.InternalId && hit.InternalId == m.TurnOfPlay {
				m.pull_card_for_player(player)
			}
		}
	}
}

func (m *SPGameRoom) OnStand(c interface{}, stand *netw.Stand) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {
		for _, player := range m.GamePlayers {
			if player.InternalId == stand.InternalId && stand.InternalId == m.TurnOfPlay {
				m.skip_next_player()
			}
		}
	}
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
func (m *SPGameRoom) skip_next_player() {
	m.CurrentPlayerCursor--
	m.next_play()
}
func (m *SPGameRoom) next_play() {
	if m.CurrentPlayerCursor > -1 {
		m.in_play_step()
	} else {
		//TODO: if system need to pull card implement
		//m.pull_card_for_system()
		m.GameStateEvent <- gs.DONE
	}
}
func (m *SPGameRoom) in_play_step() {
	player := m.GamePlayers[m.CurrentPlayerCursor]
	m.TurnOfPlay = player.InternalId
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			InternalId: player.InternalId,
			Code:       "turn_play",
		},
		MessageCode: netw.EEvent,
	}
}

func (m *SPGameRoom) pull_card_for_system() {
	card := m.PopCard()
	m.System.HitCard(card)
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Hit{
			InternalId: m.System.InternalId,
			Card:       card.String(),
			Visible:    m.System.CardVisibility(),
		},
		MessageCode: netw.EHit,
	}
}

func (m *SPGameRoom) pull_card_for_player(player *netsp.SPPlayer) {
	card := m.PopCard()
	player.HitCard(card)
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Hit{
			InternalId: player.InternalId,
			Card:       card.String(),
			Visible:    player.CardVisibility(),
		},
		MessageCode: netw.EHit,
	}
}
func (m *SPGameRoom) gameStateChanged(state gs.GameStatu) {
	switch state {
	case gs.INIT:
		m.init()
		m.GameStateEvent <- gs.WAIT_PLAYERS
	case gs.WAIT_PLAYERS:
	case gs.PREPARING:
		go m.startGame()
	case gs.IN_PLAY:
		m.in_play_step()
	case gs.DONE:
		fmt.Println("DONE")
	}
}
