package mdl

type GameSettings struct {
	Bets []*Bet `json:"bets,omitempty"`
}

type Bet struct {
	InternalId string  `json:"internal_id"`
	Amount     float32 `json:"amount"`
}
