package sp

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hayrullahcansu/mirana/core/types"
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

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	ID       int
	UserName string
	// hub      *Server
	// // The websocket connection.
	conn *websocket.Conn
	// // Buffered channel of outbound messages.
	send chan types.Envelope
	// Location Location
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
// func (c *Client) readPump() {
// 	defer func() {
// 		c.hub.unregister <- c
// 		c.conn.Close()
// 	}()
// 	c.conn.SetReadLimit(maxMessageSize)
// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.conn.SetPongHandler(func(string) error {
// 		c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 		return nil
// 	})
// 	for {
// 		_, message, err := c.conn.ReadMessage()
// 		log.Println(string(message[:]))
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		var msg json.RawMessage
// 		env := Envelope{
// 			Message: &msg,
// 		}
// 		if err := json.Unmarshal(message, &env); err != nil {
// 			log.Fatal(err)
// 		}
// 		switch env.MessageCode {
// 		case ELocation:
// 			var location Location
// 			if err := json.Unmarshal(msg, &location); err != nil {
// 				log.Fatal(err)
// 			}
// 			c.Location = location
// 		case EEvent:
// 		}
// 	}
// }

// // writePump pumps messages from the hub to the websocket connection.
// //
// // A goroutine running writePump is started for each connection. The
// // application ensures that there is at most one writer to a connection by
// // executing all writes from this goroutine.
// func (c *Client) writePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			//TextMessage denotes a text data message.
// 			// The text message payload is interpreted as UTF-8 encoded text data.
// 			w, err := c.conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			msg, _ := json.Marshal(message)
// 			log.Println(string(msg[:]))
// 			w.Write(msg)
// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

// ServeWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	IDGenarator++
	rand.Seed(time.Now().UnixNano())
	rand.Float32()
	client := &Client{
		ID:       IDGenarator,
		UserName: " ",
		// hub:      hub,
		conn: conn,
		send: make(chan types.Envelope, 1),
	}
	Manager().RequestPlayGame(client)
	// hub.register <- client
	// msg := types.Envelope{
	// 	Client:      "Server",
	// 	MessageCode: types.EEvent,
	// 	Message: types.Event{
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
