package types

type Card struct {
	CardType  CardType
	CardValue CardValue
}

func NewCardString(input string) *Card {
	return parseCard(input)
}
func NewCardData(cardType CardType, cardValue CardValue) *Card {
	return &Card{
		CardType:  cardType,
		CardValue: cardValue,
	}
}

type CardType int

const (
	_Clubs CardType = iota + 0
	_Diamonds
	_Hearts
	_Spades
)

var CardTypes = []CardType{
	_Clubs,
	_Diamonds,
	_Hearts,
	_Spades,
}

type CardValue int

const (
	_1 CardValue = iota + 0
	_2
	_3
	_4
	_5
	_6
	_7
	_8
	_9
	_10
	_JACK
	_QUEEN
	_KING
)

var CardValues = []CardValue{
	_1,
	_2,
	_3,
	_4,
	_5,
	_6,
	_7,
	_8,
	_9,
	_10,
	_JACK,
	_QUEEN,
	_KING,
}

func parseCard(input string) *Card {
	card := &Card{}
	if len(input) == 2 {
		card.CardValue = parseValue(string(input[0:1]))
		card.CardType = parseType(string(input[1:2]))
		return card
	}
	if len(input) == 3 {
		card.CardValue = parseValue(string(input[0:1]))
		card.CardType = parseType(string(input[1:2]))
		return card
	} else {
		return nil
	}
}
func parseValue(input string) CardValue {
	switch input {
	case "1":
		return _1
	case "2":
		return _2
	case "3":
		return _3
	case "4":
		return _4
	case "5":
		return _5
	case "6":
		return _6
	case "7":
		return _7
	case "8":
		return _8
	case "9":
		return _9
	case "10":
		return _10
	case "J":
		return _JACK
	case "Q":
		return _QUEEN
	case "K":
		return _KING
	default:
		return _KING

	}
}

func parseType(input string) CardType {
	switch input {
	case "C":
		return _Clubs
	case "D":
		return _Diamonds
	case "H":
		return _Hearts
	case "S":
		return _Spades
	default:
		return _Spades
	}
}
