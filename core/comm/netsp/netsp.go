package netsp

import (
	"bitbucket.org/digitdreamteam/mirana/core/comm/netw"
	"bitbucket.org/digitdreamteam/mirana/core/mdl"
	"bitbucket.org/digitdreamteam/mirana/core/types/gr"
)

type NetSPClient struct {
	*netw.BaseClient
	Players map[string]*SPPlayer
	IsDeal  bool
}

func NewClient() *NetSPClient {

	client := &NetSPClient{
		Players: make(map[string]*SPPlayer),
	}
	base := netw.NewBaseClient(client)
	client.BaseClient = base
	return client
}

type SPPlayer struct {
	Amount      float32
	InternalId  string
	Cards       []*mdl.Card
	IsSystem    bool
	Point       int
	Point2      int
	IsSplit     bool
	IsInsurance bool
	GameResult  gr.GameResult
}

func (c *NetSPClient) AddMoney(internalId string, amount float32) {
	p, ok := c.Players[internalId]
	if ok {
		p.addMoney(amount)
	} else {
		p = NewSPPlayer(internalId)
		p.addMoney(amount)
		c.Players[internalId] = p
	}
}

func (c *NetSPClient) SetInsurance(internalId string, insurance bool) {
	p, ok := c.Players[internalId]
	if ok {
		p.IsInsurance = insurance
	}
}

func (c *NetSPClient) Deal() {
	c.IsDeal = true
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

func (c *SPPlayer) addMoney(amount float32) {
	c.Amount += amount
}

func (c *SPPlayer) HitCard(card *mdl.Card) {
	c.Cards = append(c.Cards, card)
	s1 := 0
	s2 := 0
	var asExists = false
	for _, card := range c.Cards {
		if card.CardValue == mdl.CV_1 {
			asExists = true
		}
		val := card.CardValue.Value()
		s1 += val
		if asExists {
			s1 += 10
		}
	}
	if asExists {
		s2 = s1
	}
	c.Point = s1
	c.Point2 = s2
	if c.isOver21Limit() {
		c.GameResult = gr.LOSE
	}
}

func (c *SPPlayer) isOver21Limit() bool {
	if c.Point > 21 && c.Point2 > 0 {
		c.Point = c.Point2
	} else if c.Point2 > 21 {
		c.Point = c.Point
	}
	if (c.Point == 21 || c.Point2 == 21) && !c.IsSplit {
		c.GameResult = gr.BLACKJACK
	}
	if c.Point > 21 && c.Point2 > 21 {
		return false
	} else {
		return true
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
