package blackjack

type GameType int

const (
	SINGLE_DECK GameType = iota + 0
	AMERICAN
)

var GameTypes = []GameType{
	SINGLE_DECK,
	AMERICAN,
}
