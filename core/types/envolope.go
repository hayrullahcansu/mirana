package types

type Envelope struct {
	Client      string      `json:"client"`
	MessageCode MessageCode `json:"msg_code"`
	Message     interface{} `json:"msg,omitempty"`
}

type Location struct {
	Id string  `json:"id"`
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
}

type Event struct {
	Id      string     `json:"id"`
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Player  Location   `json:"player"`
	Players []Location `json:"players"`
}

// MessageCode is enumarete all message types
type MessageCode int

// MessageCode is enumarete all message types
const (
	ELocation MessageCode = iota + 0
	EEvent                // 1
)
