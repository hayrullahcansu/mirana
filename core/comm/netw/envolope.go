package netw

type Envelope struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     interface{} `json:"msg,omitempty"`
}

type Event struct {
	Id      string `json:"id"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Stamp struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
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
}

type Stand struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type Hit struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type Double struct {
	Id         string `json:"id"`
	InternalId string `json:"internal_id"`
}

type PlayGame struct {
	Id   string `json:"id"`
	Mode string `json:"mode"`
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
}
