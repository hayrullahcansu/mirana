package mdl

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
	CT_Clubs CardType = iota + 0
	CT_Diamonds
	CT_Hearts
	CT_Spades
)

var CardTypes = []CardType{
	CT_Clubs,
	CT_Diamonds,
	CT_Hearts,
	CT_Spades,
}

type CardValue int

const (
	CV_1 CardValue = iota + 0
	CV_2
	CV_3
	CV_4
	CV_5
	CV_6
	CV_7
	CV_8
	CV_9
	CV_10
	CV_JACK
	CV_QUEEN
	CV_KING
)

var CardValues = []CardValue{
	CV_1,
	CV_2,
	CV_3,
	CV_4,
	CV_5,
	CV_6,
	CV_7,
	CV_8,
	CV_9,
	CV_10,
	CV_JACK,
	CV_QUEEN,
	CV_KING,
}

func (c *Card) String() string {

	return c.CardValue.string() + c.CardType.string()
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
		return CV_1
	case "2":
		return CV_2
	case "3":
		return CV_3
	case "4":
		return CV_4
	case "5":
		return CV_5
	case "6":
		return CV_6
	case "7":
		return CV_7
	case "8":
		return CV_8
	case "9":
		return CV_9
	case "10":
		return CV_10
	case "J":
		return CV_JACK
	case "Q":
		return CV_QUEEN
	case "K":
		return CV_KING
	default:
		return CV_KING

	}
}

func parseType(input string) CardType {
	switch input {
	case "C":
		return CT_Clubs
	case "D":
		return CT_Diamonds
	case "H":
		return CT_Hearts
	case "S":
		return CT_Spades
	default:
		return CT_Spades
	}
}
func (c CardValue) string() string {
	switch c {
	case CV_1:
		return "1"
	case CV_2:
		return "2"
	case CV_3:
		return "3"
	case CV_4:
		return "4"
	case CV_5:
		return "5"
	case CV_6:
		return "6"
	case CV_7:
		return "7"
	case CV_8:
		return "8"
	case CV_9:
		return "9"
	case CV_10:
		return "10"
	case CV_JACK:
		return "J"
	case CV_QUEEN:
		return "Q"
	case CV_KING:
		return "K"
	default:
		return "K"

	}
}

func (c CardValue) Value() int {
	switch c {
	case CV_1:
		return 1
	case CV_2:
		return 2
	case CV_3:
		return 3
	case CV_4:
		return 4
	case CV_5:
		return 5
	case CV_6:
		return 6
	case CV_7:
		return 7
	case CV_8:
		return 8
	case CV_9:
		return 9
	case CV_10:
		return 10
	case CV_JACK:
		return 10
	case CV_QUEEN:
		return 10
	case CV_KING:
		return 10
	default:
		return 10

	}
}

func (c CardType) string() string {
	switch c {
	case CT_Clubs:
		return "C"
	case CT_Diamonds:
		return "D"
	case CT_Hearts:
		return "H"
	case CT_Spades:
		return "S"
	default:
		return "S"
	}
}
