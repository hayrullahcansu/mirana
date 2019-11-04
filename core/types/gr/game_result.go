package gr

type GameResult int

const (
	PLAYING GameResult = iota + 0
	WIN
	LOSE
)

var GameResults = []GameResult{
	PLAYING,
	WIN,
	LOSE,
}
