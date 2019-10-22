package types

type Player struct {
	InternalID int
	Client     *Client
	Unregister chan *Player
}

type Notify struct {
	SentBy  *Player
	Message *Envelope
}
