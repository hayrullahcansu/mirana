package netsp

import (
	"github.com/hayrullahcansu/mirana/core/comm/netw"
	"github.com/hayrullahcansu/mirana/core/mdl"
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
	Amount     float32
	InternalId string
	Cards      map[string]*mdl.Card
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

func (c *NetSPClient) Deal() {
	c.IsDeal = true
}
func NewSPPlayer(internalId string) *SPPlayer {
	return &SPPlayer{
		Amount:     0,
		InternalId: internalId,
		Cards:      make(map[string]*mdl.Card),
	}
}

func (c *SPPlayer) addMoney(amount float32) {
	c.Amount += amount
}

func (c *SPPlayer) HitCard(card *mdl.Card) {
	c.Cards[card.String()] = card
}
