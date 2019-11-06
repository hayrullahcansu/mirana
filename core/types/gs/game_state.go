package gs

var GameStatus = []GameStatu{
	INIT,
	WAIT_PLAYERS,
	PREPARING,
	ASK_INSURANCE,
	IN_PLAY,
	DONE,
}

type GameStatu int

const (
	INIT GameStatu = iota + 0
	WAIT_PLAYERS
	PREPARING
	ASK_INSURANCE
	IN_PLAY
	DONE
)

type GameState struct {
	gameStatu GameStatu
	IsDouble  bool
}

func NewGameState() *GameState {
	return &GameState{
		gameStatu: INIT,
		IsDouble:  false,
	}
}

func (gs GameState) GetNextStatu() (GameStatu, GameStatu) {
	oldStatu := gs.gameStatu
	newStatu := int(oldStatu)
	newStatu++
	newStatu = (newStatu % len(GameStatus))
	return oldStatu, GameStatu(newStatu)
}
