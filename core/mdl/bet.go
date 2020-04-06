package mdl

func NewBet(amount float32, internalId string) *Bet {
	return &Bet{
		Amount:     amount,
		InternalId: internalId,
	}
}
