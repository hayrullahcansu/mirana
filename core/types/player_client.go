package types

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var IDGenarator = 0

// Client is a middleman between the websocket Connection and the hub.
type Client struct {
	ID       int
	UserName string
	// hub      *Server
	// // The websocket Connection.
	Conn *websocket.Conn
	// // Buffered channel of outbound messages.
	Send chan *Envelope
	// Location Location
	Notify chan *Notify
	Player *Player
}

// // readPump pumps messages from the websocket Connection to the hub.

// // The application runs readPump in a per-Connection goroutine. The application
// // ensures that there is at most one reader on a Connection by executing all
// // reads from this goroutine.
func (c *Client) ReadPump() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
		c.Player.Unregister <- c.Player
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.Conn.ReadMessage()
		log.Println(string(message[:]))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		var msg json.RawMessage
		env := Envelope{
			Message: &msg,
		}
		if err := json.Unmarshal(message, &env); err != nil {
			log.Fatal(err)
		}
		if c.Notify != nil {
			switch env.MessageCode {
			case EEvent:
				var event Event
				if err := json.Unmarshal(msg, &event); err != nil {
					log.Fatal(err)
				}
				env.Message = event

			case EStamp:
				var stamp Stamp
				if err := json.Unmarshal(msg, &stamp); err != nil {
					log.Fatal(err)
				}
				env.Message = stamp

			case EAddMoney:
				var addMoney AddMoney
				if err := json.Unmarshal(msg, &addMoney); err != nil {
					log.Fatal(err)
				}
				env.Message = addMoney

			case EDeal:
				var deal Deal
				if err := json.Unmarshal(msg, &deal); err != nil {
					log.Fatal(err)
				}
				env.Message = deal

			case EStand:
				var stand Stand
				if err := json.Unmarshal(msg, &stand); err != nil {
					log.Fatal(err)
				}
				env.Message = stand

			case EHit:
				var hit Hit
				if err := json.Unmarshal(msg, &hit); err != nil {
					log.Fatal(err)
				}
				env.Message = hit

			case EDouble:
				var double Double
				if err := json.Unmarshal(msg, &double); err != nil {
					log.Fatal(err)
				}
				env.Message = double
			}
			c.Notify <- &Notify{
				Message: &env,
				SentBy:  c.Player,
			}
		}

	}
}

// // writePump pumps messages from the hub to the websocket Connection.
// //
// // A goroutine running writePump is started for each Connection. The
// // application ensures that there is at most one writer to a Connection by
// // executing all writes from this goroutine.
func (c *Client) WritePump() {
	// ticker := time.NewTicker(pingPeriod)
	// defer func() {
	// 	ticker.Stop()
	// 	c.Conn.Close()
	// }()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//TextMessage denotes a text data message.
			// The text message payload is interpreted as UTF-8 encoded text data.
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			msg, _ := json.Marshal(message)
			log.Println(string(msg[:]))
			w.Write(msg)
			if err := w.Close(); err != nil {
				return
			}
			// case <-ticker.C:
			// 	c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			// 	if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			// 		return
			// 	}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) *Client {
	Conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	IDGenarator++
	rand.Seed(time.Now().UnixNano())
	rand.Float32()

	client := &Client{
		ID:       IDGenarator,
		UserName: " ",
		// hub:      hub,
		Conn: Conn,
		Send: make(chan *Envelope, 1),
	}

	return client
	// hub.register <- client
	// msg := Envelope{
	// 	Client:      "Server",
	// 	MessageCode: EEvent,
	// 	Message: Event{
	// 		Id:      strconv.Itoa(IDGenarator),
	// 		Code:    "10",
	// 		Message: "Wellcome to room",
	// 		// Players: hub.GetUserList(),
	// 	}}
	// client.send <- msg

	// go client.writePump()
	// client.readPump()
}

// // formatRequest generates ascii representation of a request
// func formatRequest(r *http.Request) string {
// 	// Create return string
// 	var request []string
// 	// Add the request string
// 	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
// 	request = append(request, url)
// 	// Add the host
// 	request = append(request, fmt.Sprintf("Host: %v", r.Host))
// 	// Loop through headers
// 	for name, headers := range r.Header {
// 		name = strings.ToLower(name)
// 		for _, h := range headers {
// 			request = append(request, fmt.Sprintf("%v: %v", name, h))
// 		}
// 	}

// 	// If this is a POST, add post data
// 	if r.Method == "POST" {
// 		r.ParseForm()
// 		request = append(request, "\n")
// 		request = append(request, r.Form.Encode())
// 	}
// 	// Return the request as a string
// 	return strings.Join(request, "\n")
// }

// func random(min, max float32) float32 {
// 	rand.Seed(time.Now().Unix())
// 	return min + rand.Float32()*(max-min)
// }
