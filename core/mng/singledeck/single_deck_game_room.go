package singledeck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"

	"bitbucket.org/digitdreamteam/mirana/core/comm/netsp"
	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
	"bitbucket.org/digitdreamteam/mirana/core/mdl"
	"bitbucket.org/digitdreamteam/mirana/core/types/gr"
	"bitbucket.org/digitdreamteam/mirana/core/types/gs"
	"bitbucket.org/digitdreamteam/mirana/utils/que"
)

type SingleDeckGameRoom struct {
	*netw.BaseRoomManager
	PlayerConnection    *netsp.NetSPClient
	GameState           *gs.GameState
	GamePlayers         []*netsp.SPPlayer
	System              *netsp.SPPlayer
	GameStateEvent      chan gs.GameStatu
	CurrentPlayerCursor int
	TurnOfPlay          string
	Pack                *que.Queue
}

func NewSingleDeckGameRoom() *SingleDeckGameRoom {
	gameRoom := &SingleDeckGameRoom{
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
		GamePlayers:     make([]*netsp.SPPlayer, 0, 12),
		GameStateEvent:  make(chan gs.GameStatu, 1),
	}
	go gameRoom.ListenEvents()
	gameRoom.GameStateEvent <- gs.INIT
	return gameRoom
}

func (s *SingleDeckGameRoom) ListenEvents() {
	fmt.Println("GIRDI")
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
		case e := <-s.Broadcast:
			if s.PlayerConnection != nil {
				s.PlayerConnection.Send <- e
			}
		case notify := <-s.Notify:
			s.OnNotify(notify)
		// default:
		case gameStateEvent := <-s.GameStateEvent:
			s.gameStateChanged(gameStateEvent)
		}
	}
}

func (s *SingleDeckGameRoom) OnNotify(notify *netw.Notify) {
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

func (m *SingleDeckGameRoom) ConnectGame(c *netsp.NetSPClient) {
	m.PlayerConnection = c
	c.Notify = m.Notify
	m.Register <- c
}

func (m *SingleDeckGameRoom) OnConnect(c interface{}) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {

	}
}
func (m *SingleDeckGameRoom) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
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

func (m *SingleDeckGameRoom) OnAddMoney(c interface{}, addMoney *netw.AddMoney) {
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

func (m *SingleDeckGameRoom) OnSplit(c interface{}, split *netw.Split) {
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		m.split_player(client, split.InternalId)
	}
}

func (m *SingleDeckGameRoom) OnDeal(c interface{}, deal *netw.Deal) {
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
	everyoneDealed = m.PlayerConnection.IsDeal
	if everyoneDealed {
		m.GameStateEvent <- gs.PREPARING
	}
}

func (m *SingleDeckGameRoom) prepare() {
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
		if len(m.PlayerConnection.Players) > 0 && m.PlayerConnection.IsDeal {
			for _, p := range m.PlayerConnection.Players {
				m.GamePlayers = append(m.GamePlayers, p)
				indexer++
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

	//insurance check
	if m.System.HasAceFirstCard() {
		m.ask_insurance()
	}
	//split asking if check
	m.split_asking_if_check()

	m.GameStateEvent <- gs.IN_PLAY
}

func (m *SingleDeckGameRoom) split_asking_if_check() {
	for _, player := range m.GamePlayers {
		if player.CanSplit {
			m.PlayerConnection.Send <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.Event{
					Code:       "ask_split",
					InternalId: player.InternalId,
				},
				MessageCode: netw.EEvent,
			}
		}
	}
}

func (m *SingleDeckGameRoom) ask_insurance() {
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			Code:       "ask_insurance",
			InternalId: "server",
		},
		MessageCode: netw.EEvent,
	}
}

func (m *SingleDeckGameRoom) split_player(client *netsp.NetSPClient, internalId string) {
	player, ok := client.Players[internalId]
	var splitedPlayer *netsp.SPPlayer
	if ok {
		if player.CanSplit {
			secondCard := player.RemoveCard(1)
			splitedPlayer = netsp.NewSplitedSPPlayer(player)
			splitedPlayer.HitCard(secondCard)
		}
	}
	player.IsSplit = true
	card := m.PopCard()
	player.HitCard(card)
	card = m.PopCard()
	splitedPlayer.HitCard(card)
	m.GamePlayers = append(m.GamePlayers, splitedPlayer)
	m.CurrentPlayerCursor++
	copy(m.GamePlayers[m.CurrentPlayerCursor:], m.GamePlayers[m.CurrentPlayerCursor-1:])
	m.GamePlayers[m.CurrentPlayerCursor-1] = splitedPlayer
	RefCards := player.GetCardStringCommaDelemited()
	SplitedPlayerCards := splitedPlayer.GetCardStringCommaDelemited()
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Split{
			InternalId:         splitedPlayer.InternalId,
			Amount:             splitedPlayer.Amount,
			Ref:                player.InternalId,
			RefCards:           RefCards,
			SplitedPlayerCards: SplitedPlayerCards,
		},
		MessageCode: netw.ESplit,
	}
	time.Sleep(time.Millisecond * 300)

}

func (m *SingleDeckGameRoom) OnEvent(c interface{}, event *netw.Event) {
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		if event.Code == "insurance" {
			insurance := false
			if event.Message == "true" {
				insurance = true
			}
			client.SetInsurance(event.InternalId, insurance)
			// m.Broadcast <- &netw.Envelope{
			// 	Client: "client_id",
			// 	Message: &netw.AddMoney{
			// 		InternalId: addMoney.InternalId,
			// 		Amount:     addMoney.Amount,
			// 	},
			// 	MessageCode: netw.EAddMoney,
			// }
		}
	}
}

func (m *SingleDeckGameRoom) OnHit(c interface{}, hit *netw.Hit) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {
		for _, player := range m.GamePlayers {
			if player.InternalId == hit.InternalId && hit.InternalId == m.TurnOfPlay {
				m.pull_card_for_player(player)
			}
		}
	}
}

func (m *SingleDeckGameRoom) OnStand(c interface{}, stand *netw.Stand) {
	_, ok := c.(*netsp.NetSPClient)
	if ok {
		for _, player := range m.GamePlayers {
			if player.InternalId == stand.InternalId && stand.InternalId == m.TurnOfPlay {
				m.skip_next_player()
			}
		}
	}
}

func (m *SingleDeckGameRoom) PopCard() *mdl.Card {
	element := m.Pack.Dequeue()
	if element != nil {
		return element.(*mdl.Card)
	}
	return nil
}

func (m *SingleDeckGameRoom) init() {
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
func (m *SingleDeckGameRoom) skip_next_player() {
	m.CurrentPlayerCursor--
	m.next_play()
}
func (m *SingleDeckGameRoom) next_play() {
	if m.CurrentPlayerCursor > -1 {
		m.send_turn_play_message_current_player()
	} else {
		//TODO: dealer must standon soft 17
		//m.pull_card_for_system()
		m.GameStateEvent <- gs.DONE
	}
}
func (m *SingleDeckGameRoom) send_turn_play_message_current_player() {
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

func (m *SingleDeckGameRoom) pull_card_for_system() {
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

func (m *SingleDeckGameRoom) pull_card_for_player(player *netsp.SPPlayer) {
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
	if player.GameResult == gr.LOSE {
		m.Broadcast <- &netw.Envelope{
			Client: "client_id",
			Message: &netw.Event{
				InternalId: player.InternalId,
				Code:       "player_lose",
			},
			MessageCode: netw.EEvent,
		}
		fmt.Println("lose" + player.InternalId)
		m.skip_next_player()
	}
}
func (m *SingleDeckGameRoom) gameStateChanged(state gs.GameStatu) {
	switch state {
	case gs.INIT:
		m.init()
		m.GameStateEvent <- gs.WAIT_PLAYERS
	case gs.WAIT_PLAYERS:
	case gs.PREPARING:
		go m.prepare()
	case gs.PRE_START:
	case gs.IN_PLAY:
		m.send_turn_play_message_current_player()
	case gs.DONE:
		m.checkWinLose()
	}
}

func (m *SingleDeckGameRoom) checkWinLose() {
	// players := make([]*netsp.SPPlayer, 0, 6)
	// for _, p := range m.GamePlayers {

	// }
	sort.Slice(m.GamePlayers, func(i, j int) bool {
		return m.GamePlayers[i].Point < m.GamePlayers[j].Point
	})
	sort.Slice(m.GamePlayers, func(i, j int) bool {
		return m.GamePlayers[i].Point > m.GamePlayers[j].Point && m.GamePlayers[i].Point <= 21
	})
	var winnerFlag bool = false
	for _, p := range m.GamePlayers {
		if p.Point == 21 {
			p.GameResult = gr.WIN
			winnerFlag = true
		} else if p.Point < 21 && !winnerFlag {
			p.GameResult = gr.WIN
			winnerFlag = true
		} else {
			p.GameResult = gr.LOSE
		}
	}

	for _, p := range m.GamePlayers {
		if p.GameResult == gr.WIN {
			fmt.Printf("Winner %s\n", p.InternalId)
		} else {
			fmt.Printf("Loser %s\n", p.InternalId)
		}
	}
}
