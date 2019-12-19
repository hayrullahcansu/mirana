package blackjack

type Rule struct {
	StandInSoftPoint   int
	DeckNumber         int
	CardCountInOneDeck int
	SplitLimit         int
	DoubleDownLimit    int
	AceCanSplit        bool
}

func GetRules(gameType GameType) *Rule {
	switch gameType {
	case SINGLE_DECK:
		return getSingleDeckRules()
	case AMERICAN:
		return getAmericanRules()
	default:
		return nil
	}
}

func getSingleDeckRules() *Rule {
	return &Rule{
		StandInSoftPoint:   17,
		DeckNumber:         1,
		CardCountInOneDeck: 52,
		SplitLimit:         3,
		DoubleDownLimit:    3,
		AceCanSplit:        true,
	}
}

func getAmericanRules() *Rule {
	return &Rule{
		StandInSoftPoint:   17,
		DeckNumber:         4,
		CardCountInOneDeck: 52,
		SplitLimit:         4,
		DoubleDownLimit:    4,
		AceCanSplit:        true,
	}
}

/*

	-----------  SINGLE DECK  --------------
	STAND_ON_SOFT_POINT    = 17
	DECK_NUMBER            = 1
	CARD_COUNT_IN_ONE_DECK = 52
	DOUBLE_DOWN_LIMIT      = 1
	----------------------------------------


	-----------     AMERICAN     -----------
	STAND_ON_SOFT_POINT    = 17
	DECK_NUMBER            = 4
	CARD_COUNT_IN_ONE_DECK = 52
	DOUBLE_DOWN_LIMIT      = 4
	----------------------------------------

*/
