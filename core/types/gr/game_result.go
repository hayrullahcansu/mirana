package gr

type GameResult int

const (
	PLAYING GameResult = iota + 0
	WIN
	LOSE
	BLACKJACK
	PUSH
)

var GameResults = []GameResult{
	PLAYING,
	WIN,
	LOSE,
	BLACKJACK,
	PUSH,
}
