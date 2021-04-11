package blackjack

import (
	"fmt"
	"strings"

	"github.com/hayrullahcansu/mirana/core/api"
	"github.com/hayrullahcansu/mirana/core/comm/netw"
	"github.com/hayrullahcansu/mirana/core/mdl"
	"github.com/hayrullahcansu/mirana/core/types/gr"
	"github.com/hayrullahcansu/mirana/utils/stack"
)

type BlackjackClient struct {
	*netw.BaseClient
	Players        map[string]*SPPlayer
	IsDeal         bool
	IsRebet        bool
	SessionBalance float32
	Bets           *stack.Stack
}

func NewClient(userId string) *BlackjackClient {

	client := &BlackjackClient{
		Players: make(map[string]*SPPlayer),
		Bets:    stack.Init(),
	}

	base := netw.NewBaseClient(client)
	client.BaseClient = base
	client.UserId = userId
	return client
}

type SPPlayer struct {
	Amount            float32
	InsuranceAmount   float32
	InternalId        string
	Cards             []*mdl.Card
	IsSystem          bool
	Point             int
	Point2            int
	IsSplit           bool
	IsInsurance       bool
	CanSplit          bool
	DoubleDownCounter int
	GameResult        gr.GameResult
}

func (c *BlackjackClient) PlaceBet(internalId string, amount float32) bool {
	return c.placeBet(internalId, amount, false)
}

func (c *BlackjackClient) PlaceDoubleDown(internalId string) bool {
	return c.placeDoubleDown(internalId)
}

func (c *BlackjackClient) PlaceInsuranceBet(internalId string) bool {
	return c.placeBet(internalId, 0, true)
}

func (c *BlackjackClient) queueBet(internalId string, amount float32) {
	b := mdl.NewBet(amount, internalId)
	c.Bets.Push(b)
}

func (c *BlackjackClient) DequeueBet() *mdl.Bet {
	i := c.Bets.Pop()
	if i != nil {
		b := i.(*mdl.Bet)
		player := c.Players[b.InternalId]
		player.Amount -= b.Amount
		if player.Amount <= 0 {
			delete(c.Players, b.InternalId)
		}
		c.AddMoney(-b.Amount)
		return b
	}
	return nil
}

//TODO: fix .Players to GamePlayers
func (c *BlackjackClient) placeBet(internalId string, amount float32, isInsurance bool) bool {
	if greader := api.Manager().CheckAmountGreaderThan(c.UserId, amount); !greader {
		//not enough money
		return greader
	}
	p, ok := c.Players[internalId]
	if !ok {
		p = NewSPPlayer(internalId)
		c.Players[internalId] = p
	}
	if isInsurance {
		p.IsInsurance = true
		amount = p.Amount / 2
	}
	p.placeBet(amount, isInsurance)
	cost := -1 * amount
	api.Manager().AddAmount(c.UserId, cost)
	c.SessionBalance += (cost)
	if !isInsurance {
		c.queueBet(internalId, amount)
	}
	return true
}

func (c *BlackjackClient) placeDoubleDown(internalId string) bool {
	p, ok := c.Players[internalId]
	if !ok {
		return ok
	}

	if greader := api.Manager().CheckAmountGreaderThan(c.UserId, p.Amount); !greader {
		//not enough money
		return greader
	}
	amount := p.Amount

	p.placeBet(amount, false)
	cost := -1 * amount
	api.Manager().AddAmount(c.UserId, cost)
	c.SessionBalance += (cost)
	return true
}

func (c *BlackjackClient) AddMoney(amount float32) {
	api.Manager().AddAmount(c.UserId, amount)
	c.SessionBalance += (amount)
}

func (c *BlackjackClient) Reset() {
	c.IsDeal = false
	c.IsRebet = false
}

func (c *SPPlayer) Reset() {
	c.Cards = make([]*mdl.Card, 0, 10)
	c.DoubleDownCounter = 0
}

func (c *BlackjackClient) Deal() {
	c.IsDeal = true
}
func NewSplitedSPPlayer(ref *SPPlayer) *SPPlayer {
	return &SPPlayer{
		Amount:     ref.Amount,
		InternalId: ref.InternalId + "s",
		Cards:      make([]*mdl.Card, 0, 10),
		IsSystem:   false,
		IsSplit:    true,
	}
}

func NewSPPlayer(internalId string) *SPPlayer {
	return &SPPlayer{
		Amount:     0,
		InternalId: internalId,
		Cards:      make([]*mdl.Card, 0, 10),
		IsSystem:   false,
	}
}

func NewSPSystemPlayer() *SPPlayer {
	return &SPPlayer{
		Amount:     0,
		InternalId: "server",
		Cards:      make([]*mdl.Card, 0, 10),
		IsSystem:   true,
	}
}

func (c *SPPlayer) placeBet(amount float32, isInsurance bool) {
	if isInsurance {
		c.InsuranceAmount += amount
	} else {
		c.Amount += amount
	}
}

func (c *SPPlayer) RemoveCard(index int) *mdl.Card {
	tempCard := c.Cards[index]
	c.Cards = append(c.Cards[:index], c.Cards[index+1:]...)
	c.calculateScore()
	c.setCanSplit()
	return tempCard
}
func (c *SPPlayer) HitCard(card *mdl.Card) {
	c.Cards = append(c.Cards, card)
	c.calculateScore()
	c.setCanSplit()
}

func (c *SPPlayer) GetCardStringCommaDelemited() string {
	reg := make([]string, len(c.Cards))
	for i, card := range c.Cards {
		reg[i] = card.String()
	}
	return strings.Join(reg[:], ",")
}

func (c *SPPlayer) isOver21Limit() bool {
	if c.Point2 > 0 && c.Point2 <= 21 {
		c.Point = c.Point2
	}
	// if c.Point2 <= 21 && c.Point2 > c.Point {
	// 	c.Point = c.Point2
	// }
	if c.Point == 21 || c.Point2 == 21 {
		if len(c.Cards) == 2 {
			c.GameResult = gr.BLACKJACK
		} else {
			c.GameResult = gr.WIN
		}
	}
	if c.Point > 21 {
		return true
	} else {
		return false
	}
}

func (c *SPPlayer) CardVisibility() bool {
	if c.IsSystem && len(c.Cards) == 2 {
		return false
	} else {
		return true
	}
}

func (c *SPPlayer) HasAceFirstCard() bool {
	if len(c.Cards) > 0 && c.Cards[0].CardValue == mdl.CV_1 {
		return true
	}
	return false
}

func (c *SPPlayer) IsInsuranceWorked() bool {
	if c.IsSystem && len(c.Cards) == 2 && c.Cards[0].CardValue == mdl.CV_1 && c.Point == 21 {
		return true
	}
	return false
}
func (c *SPPlayer) IsBlackjack() bool {
	if len(c.Cards) == 2 && c.Point == 21 {
		return true
	}
	return false
}
func (c *SPPlayer) setCanSplit() {
	if len(c.Cards) == 2 {
		fmt.Printf(".")
	}
	if len(c.Cards) == 2 && c.Cards[0].CardValue.Value() == c.Cards[1].CardValue.Value() {
		c.CanSplit = true
	} else {
		c.CanSplit = false
	}
}
func (c *SPPlayer) calculateScore() {
	s1 := 0
	s2 := 0
	var asExists = false
	var asUsed = false
	d := fmt.Sprintf("car number: %d", len(c.Cards))
	for _, card := range c.Cards {
		d += "counted;"
		if card == nil {
			fmt.Println("error-> " + d)
		}
		if card.CardValue == mdl.CV_1 && !asExists {
			asExists = true
		}
		val := card.CardValue.Value()
		s1 += val
		if asExists && !asUsed {
			s2 += 10
			asUsed = true
		}
	}
	if asExists {
		s2 += s1
	}
	c.Point = s1
	c.Point2 = s2
	if c.isOver21Limit() {
		c.GameResult = gr.LOSE
	}
}
