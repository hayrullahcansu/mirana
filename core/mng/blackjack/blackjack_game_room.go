package blackjack

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/hayrullahcansu/mirana/core/api"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
	"github.com/hayrullahcansu/mirana/core/mdl"
	"github.com/hayrullahcansu/mirana/core/types/gr"
	"github.com/hayrullahcansu/mirana/core/types/gs"
	"github.com/hayrullahcansu/mirana/utils"
	"github.com/hayrullahcansu/mirana/utils/que"
)

type BlackjackGameRoom struct {
	*netw.BaseRoomManager
	PlayerConnection    *BlackjackClient
	GameState           *gs.GameState
	GameStatu           gs.GameStatu
	GamePlayers         []*SPPlayer
	System              *SPPlayer
	GameStateEvent      chan gs.GameStatu
	CurrentPlayerCursor int
	TurnOfPlay          string
	Pack                *que.Queue
	GameType            GameType
	Rule                *Rule
	RuleModule          *RuleModule
}

func NewBlackjackGameRoom(gameType GameType) *BlackjackGameRoom {
	rules := GetRules(gameType)
	gameRoom := &BlackjackGameRoom{
		BaseRoomManager: netw.NewBaseRoomManager(),
		GameState:       gs.NewGameState(),
		GamePlayers:     make([]*SPPlayer, 0, 12),
		GameStateEvent:  make(chan gs.GameStatu, 1),
		GameStatu:       gs.NONE,
		GameType:        gameType,
		Rule:            rules,
		RuleModule:      NewRuleModule(),
	}
	go gameRoom.ListenEvents()
	gameRoom.GameStateEvent <- gs.INIT
	return gameRoom
}

func (s *BlackjackGameRoom) PrintRoomStatus() {

	fmt.Printf("\n\n\n")
	fmt.Printf("--------Room-------\n")
	fmt.Printf("CurrentPlayerCursor:%d\n", s.CurrentPlayerCursor)
	fmt.Printf("TurnOfPlay:%s\n", s.TurnOfPlay)
	fmt.Printf("-------------------\n")
	fmt.Printf("GameStatu %+v\n", s.GameStatu)
	fmt.Printf("Client %+v\n", s.PlayerConnection)
	fmt.Printf("-----Game Players----\n")
	for i, p := range s.GamePlayers {
		fmt.Printf("Player[%d] %+v\n", i, p)

	}
	fmt.Printf("---------------------\n")
	fmt.Printf("\n\n\n")
}

func (s *BlackjackGameRoom) ListenEvents() {
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

func (m *BlackjackGameRoom) gameStateChanged(state gs.GameStatu) {
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

func (s *BlackjackGameRoom) DoUpdate(update *netw.Update) {
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

func (s *BlackjackGameRoom) SendGameConfig() {
	s.Broadcast <- &netw.Envelope{
		Client:      "client_id",
		Message:     &netw.GameConfig{},
		MessageCode: netw.EGameConfig,
	}
}

func (s *BlackjackGameRoom) OnNotify(notify *netw.Notify) {
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

func (m *BlackjackGameRoom) ConnectGame(c *BlackjackClient) {
	m.PlayerConnection = c
	c.Notify = m.Notify
	m.Register <- c
}

func (m *BlackjackGameRoom) OnConnect(c interface{}) {
	client, ok := c.(*BlackjackClient)
	if ok {
		logrus.Infof("Player Connected UserId:%s", client.UserId)
		client.Unregister = m.Unregister
		m.Update <- &netw.Update{
			Type: "account",
			Code: client.UserId,
		}
	} else {
		logrus.Error("BlackjackClient Cast Exception")
	}
}

func (m *BlackjackGameRoom) OnDisconnect(c interface{}) {
	client, ok := c.(*BlackjackClient)
	if ok {
		logrus.Infof("Player Diconnected UserId:%s", client.UserId)
		m.GameStateEvent <- gs.PURGE
	} else {
		logrus.Error("BlackjackClient Cast Exception")
	}
}

func (m *BlackjackGameRoom) PurgeRoom() {
	m.BroadcastStop <- true
	m.ListenEventsStop <- true
	m.PlayerConnection = nil
	Manager().RemoveGameRoom(m)
}

func (m *BlackjackGameRoom) OnPlayGame(c interface{}, playGame *netw.PlayGame) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*BlackjackClient)
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

func (m *BlackjackGameRoom) OnAddMoney(c interface{}, addMoney *netw.AddMoney) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*BlackjackClient)
	if ok {
		logrus.Infof("UserId:%s Req:OnAddMoney Model:%s", client.UserId, utils.ToJson(addMoney))
		if m.GameStatu == gs.WAIT_PLAYERS {
			if client.PlaceBet(addMoney.InternalId, addMoney.Amount) {
				logrus.Infof("UserId:%s PlacedBet Amount:%f", client.UserId, addMoney.Amount)
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
				logrus.Infof("UserId:%s Amount:%f Not Enough", client.UserId, addMoney.Amount)
				//TODO: send not enough balance
			}
		} else {
			logrus.Warnf("UserId:%s Invalid Request", client.UserId)
		}
	} else {
		logrus.Error("BlackjackClient Cast Exception")
	}
}

func (m *BlackjackGameRoom) OnSplit(c interface{}, split *netw.Split) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*BlackjackClient)
	if ok {
		logrus.Infof("UserId:%s Req:OnSplit Model:%s", client.UserId, utils.ToJson(split))
		if m.GameStatu == gs.IN_PLAY {
			if m.RuleModule.CheckCanSplit(split.InternalId, m.Rule) {
				m.split_player(client, split.InternalId)
			}
		} else {
			logrus.Warnf("UserId:%s Invalid Request", client.UserId)
		}
	} else {
		logrus.Error("BlackjackClient Cast Exception")
	}

}

func (m *BlackjackGameRoom) OnDeal(c interface{}, deal *netw.Deal) {
	m.L.Lock()
	client, ok := c.(*BlackjackClient)
	//TODO: check balance and other controls
	if ok && m.GameStatu == gs.WAIT_PLAYERS {
		if !client.IsRebet && (deal.Code == "rebet" || deal.Code == "rebet_and_deal") {
			var settings mdl.GameSettings
			client.IsRebet = true
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
		if deal.Code == "deal" || deal.Code == "rebet_and_deal" {
			client.Deal()
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.Deal{
					InternalId: deal.InternalId,
					Code:       "dealed",
				},
				MessageCode: netw.EDeal,
			}
		} else {
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.Deal{
					InternalId: deal.InternalId,
					Code:       "rebet",
				},
				MessageCode: netw.EDeal,
			}
		}
	}
	everyoneDealed := true
	everyoneDealed = m.PlayerConnection.IsDeal
	m.L.Unlock()
	if everyoneDealed {
		m.GameStateEvent <- gs.PREPARING
	}
}

func (m *BlackjackGameRoom) OnEvent(c interface{}, event *netw.Event) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*BlackjackClient)
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
		}
	} else if ok && m.GameStatu == gs.WAIT_PLAYERS && event.Code == "undo_bet" {
		lastBet := client.DequeueBet()
		if lastBet != nil {
			m.Broadcast <- &netw.Envelope{
				Client: "client_id",
				Message: &netw.AddMoney{
					InternalId: lastBet.InternalId,
					Amount:     lastBet.Amount,
					Op:         "undo_bet",
				},
				MessageCode: netw.EAddMoney,
			}
		}
	}
}

func (m *BlackjackGameRoom) OnHit(c interface{}, hit *netw.Hit) {
	m.L.Lock()
	defer m.L.Unlock()
	_, ok := c.(*BlackjackClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		for _, player := range m.GamePlayers {
			if player.InternalId == hit.InternalId && hit.InternalId == m.TurnOfPlay {
				if m.pull_card_for_player(player) {
					m.skip_next_player()
				}
			}
		}
	}
}

func (m *BlackjackGameRoom) OnDouble(c interface{}, double *netw.Double) {
	m.L.Lock()
	defer m.L.Unlock()
	client, ok := c.(*BlackjackClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		for _, player := range m.GamePlayers {
			if player.InternalId == double.InternalId && double.InternalId == m.TurnOfPlay {
				if m.RuleModule.CheckCanDoubleDown(double.InternalId, m.Rule) {
					dd_ok := m.double_down_for_player(client, double.InternalId)
					if dd_ok {
						m.skip_next_player()
					}
				}
			}
		}
	}
}

func (m *BlackjackGameRoom) OnStand(c interface{}, stand *netw.Stand) {
	m.L.Lock()
	defer m.L.Unlock()
	_, ok := c.(*BlackjackClient)
	if ok && m.GameStatu == gs.IN_PLAY {
		for _, player := range m.GamePlayers {
			if player.InternalId == stand.InternalId && stand.InternalId == m.TurnOfPlay {
				m.skip_next_player()
			}
		}
	}
}

func (m *BlackjackGameRoom) PopCard() *mdl.Card {
	// return &mdl.Card{
	// 	CardType:  mdl.CT_Diamonds,
	// 	CardValue: mdl.CV_JACK,
	// }
	element := m.Pack.Dequeue()
	if element != nil {
		return element.(*mdl.Card)
	}
	return nil
}

func (m *BlackjackGameRoom) init() {
	switch m.GameType {
	case SINGLE_DECK:
		m.Pack = utils.GetSingleDeckPack()
		break
	case AMERICAN:
		m.Pack = utils.GetAmericanPack()
		break
	default:
		logrus.Errorf("Invalid Game Type")
	}
}
func (m *BlackjackGameRoom) resetGame(justClear bool) {
	if m.PlayerConnection != nil {
		m.PlayerConnection.Reset()
	}
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
func (m *BlackjackGameRoom) prepare() {
	m.L.Lock()
	logrus.Warn("deck size:%d  limit:%d", len(m.Pack.Values), m.Rule.CardCountInOneDeck*m.Rule.DeckNumber/2)
	if len(m.Pack.Values) < m.Rule.CardCountInOneDeck*m.Rule.DeckNumber/2 {
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
	// m.split_asking_if_check()
	m.L.Unlock()
	m.GameStateEvent <- gs.IN_PLAY
}

func (m *BlackjackGameRoom) split_asking_if_check() {
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

func (m *BlackjackGameRoom) ask_insurance() {
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			Code:       "ask_insurance",
			InternalId: "server",
		},
		MessageCode: netw.EEvent,
	}
}

func (m *BlackjackGameRoom) split_player(client *BlackjackClient, internalId string) {
	player, ok := client.Players[internalId]
	var splitedPlayer *SPPlayer
	if ok {
		if player.CanSplit && api.Manager().CheckAmountGreaderThan(client.UserId, player.Amount) {
			splitedPlayer = NewSplitedSPPlayer(player)
			//this line comment out incorrect split operation.
			// card = m.PopCard()
			// splitedPlayer.HitCard(card)
			m.GamePlayers = append(m.GamePlayers, splitedPlayer)
			m.CurrentPlayerCursor++
			copy(m.GamePlayers[m.CurrentPlayerCursor:], m.GamePlayers[m.CurrentPlayerCursor-1:])
			m.GamePlayers[m.CurrentPlayerCursor-1] = splitedPlayer
			// RefCards := player.GetCardStringCommaDelemited()
			// SplitedPlayerCards := splitedPlayer.GetCardStringCommaDelemited()
			if client.PlaceBet(splitedPlayer.InternalId, player.Amount) {
				m.RuleModule.IncreaseSplitCounter(internalId)
				secondCard := player.RemoveCard(1)
				splitedPlayer.HitCard(secondCard)
				player.IsSplit = true
				m.Update <- &netw.Update{
					Type: "account",
					Code: client.UserId,
				}
				m.Broadcast <- &netw.Envelope{
					Client: "client_id",
					Message: &netw.Split{
						InternalId: splitedPlayer.InternalId,
						Amount:     splitedPlayer.Amount,
						Ref:        player.InternalId,
						// RefCards:           RefCards,
						// SplitedPlayerCards: SplitedPlayerCards,
					},
					MessageCode: netw.ESplit,
				}
				if m.pull_card_for_player(player) {
					m.skip_next_player()
				}
			} else {
				//Remove splitted player
				m.GamePlayers = append(m.GamePlayers[:m.CurrentPlayerCursor-1], m.GamePlayers[m.CurrentPlayerCursor:]...)
				m.CurrentPlayerCursor--
			}

			time.Sleep(time.Millisecond * 300)
		}
	}
}

func (m *BlackjackGameRoom) skip_next_player() {
	m.CurrentPlayerCursor--
	m.next_play()
}
func (m *BlackjackGameRoom) next_play() {
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
		for m.System.Point < m.Rule.StandInSoftPoint {
			m.pull_card_for_system()
			time.Sleep(time.Second * 1)
		}
		time.Sleep(time.Millisecond * 300)
		m.GameStateEvent <- gs.DONE
	}
}
func (m *BlackjackGameRoom) send_turn_play_message_current_player() {
	player := m.GamePlayers[m.CurrentPlayerCursor]
	if player.Point >= 21 {
		m.skip_next_player()
		return
	}
	m.TurnOfPlay = player.InternalId
	if player.IsSplit {
		player.IsSplit = false
		if m.pull_card_for_playerfor(player) {
			m.skip_next_player()
		}

	}
	message := ""
	if player.CanSplit && m.RuleModule.CheckCanSplit(player.InternalId, m.Rule) {
		message += "show_split_button "
	} else {
		message += "hide_split_button "
	}
	if m.RuleModule.CheckCanDoubleDown(player.InternalId, m.Rule) {
		message += "show_double_down_button "
	} else {
		message += "hide_double_down_button "
	}

	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			InternalId: player.InternalId,
			Code:       "turn_play",
			Message:    message,
		},
		MessageCode: netw.EEvent,
	}
}

func (m *BlackjackGameRoom) pull_card_for_system() {
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

func (m *BlackjackGameRoom) pull_card_for_player(player *SPPlayer) bool {
	card := m.PopCard()
	if card == nil {
		fmt.Println(fmt.Sprintf("BUG OLUSTU CUNKU BITTI %d", len(m.Pack.Values)))
	}
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
		return true
	} else if player.Point >= 21 {
		return true
	}
	return false
}

func (m *BlackjackGameRoom) pull_card_for_playerfor(player *SPPlayer) bool {
	firstCard := player.Cards[0]
	card := mdl.NewCardData(firstCard.CardType, firstCard.CardValue)

	if card == nil {
		fmt.Println(fmt.Sprintf("BUG OLUSTU CUNKU BITTI %d", len(m.Pack.Values)))
	}
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
		return true
	} else if player.Point >= 21 {
		return true
	}
	return false
}
func (m *BlackjackGameRoom) double_down_for_player(client *BlackjackClient, internalId string) bool {
	player, ok := client.Players[internalId]
	if ok {
		if player.DoubleDownCounter < m.Rule.DoubleDownLimit && api.Manager().CheckAmountGreaderThan(client.UserId, player.Amount) {
			dd_ok := client.PlaceDoubleDown(internalId)
			if dd_ok {
				m.RuleModule.IncreaseDoubleDownCounter(internalId)
				m.Broadcast <- &netw.Envelope{
					Client: "client_id",
					Message: &netw.Double{
						InternalId: player.InternalId,
					},
					MessageCode: netw.EDouble,
				}
				m.Update <- &netw.Update{
					Type: "account",
					Code: client.UserId,
				}
				m.pull_card_for_player(player)
				return true
			}
		}
	}
	return false
}
func (m *BlackjackGameRoom) checkWinLose() {
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
			} else if p.Point == 21 {
				if p.IsBlackjack() {
					p.GameResult = gr.BLACKJACK
				} else {
					p.GameResult = gr.WIN
				}
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
			} else if p.Point == 21 {
				if p.IsBlackjack() {
					p.GameResult = gr.BLACKJACK
				} else {
					p.GameResult = gr.WIN
				}
			} else if p.Point > m.System.Point {
				p.GameResult = gr.WIN
			} else if p.Point == m.System.Point {
				p.GameResult = gr.PUSH
			} else {
				p.GameResult = gr.LOSE
			}
		}
	}
	var winMoney, total float32
	for _, p := range m.GamePlayers {
		if p.IsInsurance && m.System.IsInsuranceWorked() {
			winMoney = p.Amount
			m.PlayerConnection.AddMoney(winMoney)
			fmt.Printf("Insurance %s\n", p.InternalId)
		}
		messsage := "player_lose"
		id := p.InternalId
		switch p.GameResult {
		case gr.WIN:
			winMoney = p.Amount * 2
			m.PlayerConnection.AddMoney(winMoney)
			fmt.Printf("Winner %s\n", p.InternalId)
			messsage = "player_win"
		case gr.LOSE:
			winMoney = 0
			fmt.Printf("Loser %s\n", p.InternalId)
			messsage = "player_lose"
		case gr.PUSH:
			winMoney = p.Amount
			m.PlayerConnection.AddMoney(winMoney)
			fmt.Printf("Push %s\n", p.InternalId)
			messsage = "player_push"
		case gr.BLACKJACK:
			winMoney = p.Amount * 5 / 2
			m.PlayerConnection.AddMoney(winMoney)
			fmt.Printf("Blackjack %s\n", p.InternalId)
			messsage = "player_blackjack"
		default:
			winMoney = 0
			fmt.Printf("Default %s\n", p.InternalId)
		}
		total += winMoney
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
			Message:    fmt.Sprintf("%f", total),
		},
		MessageCode: netw.EEvent,
	}
	m.Update <- &netw.Update{
		Type: "account",
		Code: m.PlayerConnection.UserId,
	}
	time.Sleep(time.Second * 5)
	m.GameStateEvent <- gs.INIT
	m.Broadcast <- &netw.Envelope{
		Client: "client_id",
		Message: &netw.Event{
			InternalId: "server",
			Code:       "you_can_play_again",
		},
		MessageCode: netw.EEvent,
	}
	m.resetGame(true)
}
