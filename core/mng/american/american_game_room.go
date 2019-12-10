package american

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"

	"bitbucket.org/digitdreamteam/mirana/core/api"
	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
	"bitbucket.org/digitdreamteam/mirana/core/mdl"
	"bitbucket.org/digitdreamteam/mirana/core/types/gr"
	"bitbucket.org/digitdreamteam/mirana/core/types/gs"
	"bitbucket.org/digitdreamteam/mirana/utils"
	"bitbucket.org/digitdreamteam/mirana/utils/que"
)

const (
	STAND_ON_SOFT_POINT    = 17
	DECK_NUMBER            = 4
	CARD_COUNT_IN_ONE_DECK = 52
)

type AmericanGameRoom struct {
	*netw.BaseRoomManager
	PlayerConnection    *AmericanSPClient
	GameState           *gs.GameState
	GameStatu           gs.GameStatu
	GamePlayers         []*SPPlayer
	System              *SPPlayer
	GameStateEvent      chan gs.GameStatu
	CurrentPlayerCursor int
	TurnOfPlay          string
	Pack                *que.Queue
}

func NewAmericanGameRoom() *AmericanGameRoom {
	gameRoom := &AmericanGameRoom{
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
		GamePlayers:     make([]*SPPlayer, 0, 12),
		GameStateEvent:  make(chan gs.GameStatu, 1),
		GameStatu:       gs.NONE,
	}
	go gameRoom.ListenEvents()
	gameRoom.GameStateEvent <- gs.INIT
	return gameRoom
}

func (s *AmericanGameRoom) ListenEvents() {
	fmt.Println("GIRDI")
	go func() {
		for {
			select {
			case e := <-s.Broadcast:
				if s.PlayerConnection != nil {
					s.PlayerConnection.Send <- e
				}
			case e := <-s.BroadcastStop:
				if e {
					return
				}
			}
		}
	}()
	for {
		select {
		case player := <-s.Register:
			s.OnConnect(player)
		case player := <-s.Unregister:
			s.OnDisconnect(player)
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
		case e := <-s.ListenEventsStop:
			if e {
				return
			}
		}
	}
}

func (m *AmericanGameRoom) gameStateChanged(state gs.GameStatu) {
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
	case gs.PURGE:
		m.PurgeRoom()
	}
}

func (s *AmericanGameRoom) DoUpdate(update *netw.Update) {
	switch update.Type {
	case "account":
		if s.PlayerConnection.UserId == update.Code {
			u := api.Manager().GetUser(s.PlayerConnection.UserId)
			s.PlayerConnection.Send <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.User{
					UserId:     u.UserId,
					Balance:    u.Balance,
					Name:       u.Name,
					Win:        u.Win,
					Lose:       u.Lose,
					Push:       u.Push,
					Blackjack:  u.Blackjack,
					WinBalance: s.PlayerConnection.SessionBalance,
				},
				MessageCode: netw.EUser,
			}
		}
		break
	}
}

func (s *AmericanGameRoom) OnNotify(notify *netw.Notify) {
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

func (m *AmericanGameRoom) ConnectGame(c *AmericanSPClient) {
	m.PlayerConnection = c
	c.Notify = m.Notify
	m.Register <- c
}

func (m *AmericanGameRoom) OnConnect(c interface{}) {
	client, ok := c.(*AmericanSPClient)
	if ok {
		client.Unregister = m.Unregister
		m.Update <- &netw.Update{
			Type: "account",
			Code: client.UserId,
		}
	}
}
func (m *AmericanGameRoom) OnDisconnect(c interface{}) {
	_, ok := c.(*AmericanSPClient)
	if ok {
		fmt.Println("OnDisconnect")
		m.GameStateEvent <- gs.PURGE
	}
}

func (m *AmericanGameRoom) PurgeRoom() {
	m.BroadcastStop <- true
	m.ListenEventsStop <- true
	m.PlayerConnection = nil

	Manager().RemoveGameRoom(m)
}

func (m *AmericanGameRoom) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.WAIT_PLAYERS {
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

func (m *AmericanGameRoom) OnAddMoney(c interface{}, addMoney *netw.AddMoney) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.WAIT_PLAYERS {
		if client.PlaceBet(addMoney.InternalId, addMoney.Amount) {
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

func (m *AmericanGameRoom) OnSplit(c interface{}, split *netw.Split) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		m.split_player(client, split.InternalId)
	}
}

func (m *AmericanGameRoom) OnDeal(c interface{}, deal *netw.Deal) {
	m.L.Lock()
	client, ok := c.(*AmericanSPClient)
	//TODO: check balance and other controls
	if ok && m.GameStatu == gs.WAIT_PLAYERS {
		if deal.Code == "deal_new_game" && m.GameStatu == gs.DONE {
			m.resetGame(true)
			var settings mdl.GameSettings
			bytes := []byte(deal.Payload)
			if err := json.Unmarshal(bytes, &settings); err != nil {
				log.Fatal(err)
			}
			client.Players = make(map[string]*SPPlayer)
			if len(settings.Bets) > 0 {
				for _, bet := range settings.Bets {
					if bet.InternalId != "" && bet.Amount > 0 {
						if client.PlaceBet(bet.InternalId, bet.Amount) {
							m.Update <- &netw.Update{
								Type: "account",
								Code: client.UserId,
							}
							m.Broadcast <- &netw.Envelope{
								Client: "client_id",
								Message: &netw.AddMoney{
									InternalId: bet.InternalId,
									Amount:     bet.Amount,
									Op:         "set",
								},
								MessageCode: netw.EAddMoney,
							}
						}
					}
				}
			}
		}
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

func (m *AmericanGameRoom) OnEvent(c interface{}, event *netw.Event) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		if event.Code == "insurance" {
			if event.Message == "true" {
				if client.PlaceInsuranceBet(event.InternalId) {
					m.Update <- &netw.Update{
						Type: "account",
						Code: m.PlayerConnection.UserId,
					}
				}
			}
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

func (m *AmericanGameRoom) OnHit(c interface{}, hit *netw.Hit) {
	m.L.Lock()
	defer m.L.Unlock()
	_, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		for _, player := range m.GamePlayers {
			if player.InternalId == hit.InternalId && hit.InternalId == m.TurnOfPlay {
				m.pull_card_for_player(player)
			}
		}
	}
}

func (m *AmericanGameRoom) OnStand(c interface{}, stand *netw.Stand) {
	m.L.Lock()
	defer m.L.Unlock()
	_, ok := c.(*AmericanSPClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		for _, player := range m.GamePlayers {
			if player.InternalId == stand.InternalId && stand.InternalId == m.TurnOfPlay {
				m.skip_next_player()
			}
		}
	}
}

func (m *AmericanGameRoom) PopCard() *mdl.Card {
	element := m.Pack.Dequeue()
	if element != nil {
		return element.(*mdl.Card)
	}
	return nil
}

func (m *AmericanGameRoom) init() {
	m.Pack = utils.GetAmericanPack()
}
func (m *AmericanGameRoom) resetGame(justClear bool) {
	temp := make([]*SPPlayer, 0, 12)
	if !justClear {
		for _, player := range m.GamePlayers {
			if !player.IsSplit {
				player.Reset()
				temp = append(temp, player)
			}
		}
	}
	m.GamePlayers = temp
}

func (m *AmericanGameRoom) prepare() {
	m.L.Lock()
	if len(m.Pack.Values) < CARD_COUNT_IN_ONE_DECK*DECK_NUMBER/2 {
		//generate new cards
		m.init()
	}
	m.System = NewSPSystemPlayer()
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

func (m *AmericanGameRoom) split_asking_if_check() {
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

func (m *AmericanGameRoom) ask_insurance() {
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			Code:       "ask_insurance",
			InternalId: "server",
		},
		MessageCode: netw.EEvent,
	}
}

func (m *AmericanGameRoom) split_player(client *AmericanSPClient, internalId string) {
	player, ok := client.Players[internalId]
	var splitedPlayer *SPPlayer
	if ok {
		if player.CanSplit && api.Manager().CheckAmountGreaderThan(client.UserId, player.Amount) {
			secondCard := player.RemoveCard(1)
			splitedPlayer = NewSplitedSPPlayer(player)
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
			if client.PlaceBet(splitedPlayer.InternalId, player.Amount) {
				m.Update <- &netw.Update{
					Type: "account",
					Code: client.UserId,
				}
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
			} else {

			}

			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (m *AmericanGameRoom) skip_next_player() {
	m.CurrentPlayerCursor--
	m.next_play()
}
func (m *AmericanGameRoom) next_play() {
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
func (m *AmericanGameRoom) send_turn_play_message_current_player() {
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

func (m *AmericanGameRoom) pull_card_for_system() {
	card := m.PopCard()
	//TODO: remove it, because i added for test insurance
	// if len(m.System.Cards) == 0 {
	// 	card.CardValue = mdl.CV_1
	// } else if len(m.System.Cards) == 1 {
	// 	card.CardValue = mdl.CV_JACK
	// }
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

func (m *AmericanGameRoom) pull_card_for_player(player *SPPlayer) {
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
		m.skip_next_player()
	}
}

func (m *AmericanGameRoom) checkWinLose() {
	// players := make([]*SPPlayer, 0, 6)
	// for _, p := range m.GamePlayers {

	// }
	sort.Slice(m.GamePlayers, func(i, j int) bool {
		return m.GamePlayers[i].Point < m.GamePlayers[j].Point
	})
	sort.Slice(m.GamePlayers, func(i, j int) bool {
		return m.GamePlayers[i].Point > m.GamePlayers[j].Point && m.GamePlayers[i].Point <= 21
	})
	// var winnerFlag bool = false
	fmt.Printf("element size:%d\n", len(m.GamePlayers))
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
			} else {
				p.GameResult = gr.LOSE
			}
		}
	}
	for _, p := range m.GamePlayers {
		if p.IsInsurance && m.System.IsInsuranceWorked() {
			fmt.Printf("Insurance %s\n", p.InternalId)
			m.PlayerConnection.AddMoney(p.Amount)
		}
		messsage := "player_lose"
		id := p.InternalId
		switch p.GameResult {
		case gr.WIN:
			m.PlayerConnection.AddMoney(p.Amount * 2)
			fmt.Printf("Winner %s\n", p.InternalId)
			messsage = "player_win"
		case gr.LOSE:
			fmt.Printf("Loser %s\n", p.InternalId)
			messsage = "player_lose"
		case gr.PUSH:
			m.PlayerConnection.AddMoney(p.Amount)
			fmt.Printf("Push %s\n", p.InternalId)
			messsage = "player_push"
		case gr.BLACKJACK:
			m.PlayerConnection.AddMoney(p.Amount * 5 / 2)
			fmt.Printf("Blackjack %s\n", p.InternalId)
			messsage = "player_blackjack"
		default:
			fmt.Printf("Default %s\n", p.InternalId)
		}
		m.Broadcast <- &netw.Envelope{
			Client: "client_id",
			Message: &netw.Event{
				InternalId: id,
				Code:       messsage,
			},
			MessageCode: netw.EEvent,
		}
	}
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			InternalId: "server",
			Code:       "game_done",
		},
		MessageCode: netw.EEvent,
	}
	m.Update <- &netw.Update{
		Type: "account",
		Code: m.PlayerConnection.UserId,
	}

}
