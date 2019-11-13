package singledeck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/google/uuid"

	"bitbucket.org/digitdreamteam/mirana/core/api"
	"bitbucket.org/digitdreamteam/mirana/core/comm/netsp"
	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
	"bitbucket.org/digitdreamteam/mirana/core/mdl"
	"bitbucket.org/digitdreamteam/mirana/core/types/gr"
	"bitbucket.org/digitdreamteam/mirana/core/types/gs"
	"bitbucket.org/digitdreamteam/mirana/utils/que"
)

const (
	STAND_ON_SOFT_POINT = 17
)

type SingleDeckGameRoom struct {
	*netw.BaseRoomManager
	PlayerConnection    *netsp.NetSPClient
	GameState           *gs.GameState
	GameStatu           gs.GameStatu
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
		GameStatu:       gs.NONE,
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
			go func() {
				if s.PlayerConnection != nil {
					s.PlayerConnection.Send <- e
				}
			}()
		case u := <-s.Update:
			go func() {
				s.DoUpdate(u)
			}()
		case notify := <-s.Notify:
			go func() {
				s.OnNotify(notify)
			}()
		// default:
		case gameStateEvent := <-s.GameStateEvent:
			go func() {
				s.gameStateChanged(gameStateEvent)
			}()
		}
	}
}

func (s *SingleDeckGameRoom) DoUpdate(update *netw.Update) {
	switch update.Type {
	case "account":
		if s.PlayerConnection.UserId == update.Code {
			u := api.Manager().GetUser(s.PlayerConnection.UserId)
			s.PlayerConnection.Send <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.User{
					UserId:    u.UserId,
					Balance:   u.Balance,
					Name:      u.Name,
					Win:       u.Win,
					Lose:      u.Lose,
					Push:      u.Push,
					Blackjack: u.Blackjack,
				},
				MessageCode: netw.EUser,
			}
		}
		break
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
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		m.Update <- &netw.Update{
			Type: "account",
			Code: client.UserId,
		}
	}
}
func (m *SingleDeckGameRoom) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	m.L.Lock()
	defer m.L.Unlock()
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
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		if client.AddMoney(addMoney.InternalId, addMoney.Amount) {
			m.Update <- &netw.Update{
				Type: "account",
				Code: client.UserId,
			}
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.AddMoney{
					InternalId: addMoney.InternalId,
					Amount:     addMoney.Amount,
				},
				MessageCode: netw.EAddMoney,
			}
		} else {
			//TODO: send not enough balance
		}
	}
}

func (m *SingleDeckGameRoom) OnSplit(c interface{}, split *netw.Split) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*netsp.NetSPClient)
	if ok {
		m.split_player(client, split.InternalId)
	}
}

func (m *SingleDeckGameRoom) OnDeal(c interface{}, deal *netw.Deal) {
	m.L.Lock()
	client, ok := c.(*netsp.NetSPClient)
	//TODO: check balance and other controls
	if m.GameStatu == gs.DONE {
		m.resetGame()
	}
	if ok {
		client.Deal()
		m.Broadcast <- &netw.Envelope{
			Client: "client_id",
			Message: &netw.Deal{
				InternalId: deal.InternalId,
				Code:       "dealed",
			},
			MessageCode: netw.EDeal,
		}
	}
	everyoneDealed := true
	everyoneDealed = m.PlayerConnection.IsDeal
	m.L.Unlock()
	if everyoneDealed {
		m.GameStateEvent <- gs.PREPARING
	}
}

func (m *SingleDeckGameRoom) resetGame() {
	m.L.Lock()
	temp := make([]*netsp.SPPlayer, 0, 12)
	for _, player := range m.GamePlayers {
		if !player.IsSplit {
			player.Reset()
			temp = append(temp, player)
		}
	}
	m.GamePlayers = temp
	m.L.Unlock()
	m.GameStateEvent <- gs.INIT
}

func (m *SingleDeckGameRoom) prepare() {
	m.L.Lock()
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
	m.L.Unlock()
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
	}
}

func (m *SingleDeckGameRoom) OnEvent(c interface{}, event *netw.Event) {
	m.L.Lock()
	defer m.L.Unlock()
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
	m.L.Lock()
	defer m.L.Unlock()
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
	m.L.Lock()
	defer m.L.Unlock()
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
		go func() {
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.Event{
					InternalId: m.System.InternalId,
					Code:       "show_second_card",
				},
				MessageCode: netw.EEvent,
			}
		}()
		time.Sleep(time.Second * 2)
		for m.System.Point < STAND_ON_SOFT_POINT {
			m.pull_card_for_system()
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Millisecond * 300)
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
		fmt.Printf("Loser %s\n", player.InternalId)
		m.skip_next_player()
	}
}
func (m *SingleDeckGameRoom) gameStateChanged(state gs.GameStatu) {
	m.GameStatu = state
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
	// var winnerFlag bool = false

	for _, p := range m.GamePlayers {
		if m.System.Point > 21 {
			if p.Point < 21 {
				p.GameResult = gr.WIN
			} else {
				p.GameResult = gr.LOSE
			}
		} else if m.System.Point == 21 {
			if p.Point == 21 {
				p.GameResult = gr.PUSH
			} else {
				p.GameResult = gr.LOSE
			}
		} else {
			// equals m.System.Point < 21
			if p.Point > 21 {
				p.GameResult = gr.LOSE
			} else if p.Point == 21 || p.Point > m.System.Point {
				p.GameResult = gr.WIN
			} else if p.Point == m.System.Point {
				p.GameResult = gr.PUSH
			}
		}
	}
	for _, p := range m.GamePlayers {
		switch p.GameResult {
		case gr.WIN:
			fmt.Printf("Winner %s\n", p.InternalId)
		case gr.LOSE:
			fmt.Printf("Loser %s\n", p.InternalId)
		case gr.PUSH:
			fmt.Printf("Push %s\n", p.InternalId)
		case gr.BLACKJACK:
			fmt.Printf("Blackjack %s\n", p.InternalId)
		default:

		}
	}
}
