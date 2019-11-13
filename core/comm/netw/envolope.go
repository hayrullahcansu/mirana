package netw

type Envelope struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     interface{} `json:"msg,omitempty"`
}

type EnvelopeStaging struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     string      `json:"msg,omitempty"`
}

type Event struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

type Stamp struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type Split struct {
	Id                 string  `json:"id"`
	InternalId         string  `json:"internal_id"`
	Ref                string  `json:"ref"`
	Amount             float32 `json:"amount"`
	RefCards           string  `json:"ref_cards"`
	SplitedPlayerCards string  `json:"splitted_cards"`
}

type AddMoney struct {
	Id         string  `json:"id"`
	InternalId string  `json:"internal_id"`
	Amount     float32 `json:"amount"`
	Op         string  `json:"op"`
}

type Deal struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
	Code       string `json:"code"`
}

type Stand struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type Hit struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
	Card       string `json:"card"`
	Visible    bool   `json:"visible"`
}

type Double struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type PlayGame struct {
	Id   string `json:"id"`
	Mode string `json:"mode"`
}

type User struct {
	UserId     string  `json:"user_id"`
	Name       string  `json:"name"`
	Balance    float32 `json:"balance"`
	WinBalance float32 `json:"win_balance"`
	Win        int     `json:"win"`
	Lose       int     `json:"lose"`
	Push       int     `json:"push"`
	Blackjack  int     `json:"blackjack"`
}

// MessageCode is enumarete all message types
type MessageCode int

// MessageCode is enumarete all message types
const (
	EEvent    MessageCode = iota + 0
	EStamp                // 1
	EAddMoney             // 2
	EDeal                 // 3
	EStand                // 4
	EHit                  // 5
	EDouble               // 6
	EPlayGame             // 7
	ESplit                // 8
	EUser                 // 9
)

var MessageCodes = []MessageCode{
	EEvent,
	EStamp,
	EAddMoney,
	EDeal,
	EStand,
	EHit,
	EDouble,
	EPlayGame,
	ESplit,
	EUser,
}
