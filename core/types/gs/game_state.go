package gs

var GameStatus = []GameStatu{
	NONE,
	INIT,
	WAIT_PLAYERS,
	PREPARING,
	PRE_START,
	IN_PLAY,
	DONE,
}

type GameStatu int

const (
	NONE GameStatu = iota + 0
	INIT
	WAIT_PLAYERS
	PREPARING
	PRE_START
	IN_PLAY
	DONE
	PURGE
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
