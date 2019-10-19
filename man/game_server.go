package man

// import (
// 	"strconv"
// 	"sync"
// 	"time"
// )

// type Server struct {
// 	clients    map[*Client]bool
// 	broadcast  chan Envelope
// 	register   chan *Client
// 	unregister chan *Client
// 	gamestate  chan bool
// 	m          sync.Mutex
// }

// func NewServer() *Server {
// 	return &Server{
// 		broadcast:  make(chan Envelope),
// 		register:   make(chan *Client),
// 		unregister: make(chan *Client),
// 		clients:    make(map[*Client]bool),
// 		gamestate:  make(chan bool),
// 	}
// }
// func (h *Server) Run() {

// 	go func() {
// 		for {
// 			h.gamestate <- true
// 			time.Sleep(20 * time.Millisecond)
// 		}
// 	}()
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.clients[client] = true
// 			h.NewPlayerJoined(client)
// 		case client := <-h.unregister:
// 			if _, ok := h.clients[client]; ok {
// 				h.LeavedPlayer(client)
// 				delete(h.clients, client)
// 				close(client.send)
// 			}
// 		case message := <-h.broadcast:
// 			for client := range h.clients {
// 				select {
// 				case client.send <- message:
// 				default:
// 					close(client.send)
// 					delete(h.clients, client)
// 				}
// 			}
// 		case <-h.gamestate:
// 			h.SendGameState()
// 		default:
// 		}
// 	}
// }
// func (s *Server) SpecialRequestFromClient(c *Client, m *Envelope) {
// 	//content code equals is 20-29 means
// 	// about special request from client
// 	// msg := Message{Client: "Server", ContentCode: 20}
// 	// switch m.Content {
// 	// case "-help":
// 	// 	msg.Content = s.GetHelp()
// 	// 	c.send <- msg
// 	// case "-list":
// 	// 	msg.Content = s.GetUserList()
// 	// 	c.send <- msg
// 	// }
// }
// func (s *Server) SpecialRequestFromServer(c *Client, m *Envelope) {
// 	//content code equals is 30-39 means
// 	// about special request from server
// 	// if m.ContentCode == 31 {
// 	// 	switch m.Content {
// 	// 	case "UserName":
// 	// 		//save usernick to map of clients
// 	// 		c.UserName = m.Client
// 	// 		msg := Message{Client: "Server", ContentCode: 1}
// 	// 		msg.Content = m.Client + ", Wellcome to our Chatroom. Have a nice chats!"
// 	// 		c.send <- msg
// 	// 	}
// 	// }
// }
// func (s *Server) GetHelp() string {
// 	h := "Writable special request codes are below\n\n" +
// 		"-help     helps\n" +
// 		"-list     gets online user lists\n"
// 	return h
// }
// func (s *Server) GetUserList() []Location {
// 	// defer s.m.Lock()
// 	// s.m.Lock()
// 	v := make([]Location, len(s.clients))
// 	i := 0
// 	for client, _ := range s.clients {
// 		v[i] = Location{
// 			Id: strconv.Itoa(client.ID),
// 			X:  client.Location.X,
// 			Y:  client.Location.Y,
// 		}
// 		i++
// 	}
// 	return v
// }
// func (s *Server) SendGameState() {
// 	for client := range s.clients {
// 		for c := range s.clients {
// 			message := Envelope{
// 				Client:      strconv.Itoa(c.ID),
// 				MessageCode: ELocation,
// 				Message:     c.Location,
// 			}
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				delete(s.clients, client)
// 			}
// 		}
// 	}
// }

// func (s *Server) NewPlayerJoined(c *Client) {
// 	e := Event{
// 		Code:   "11",
// 		Id:     strconv.Itoa(c.ID),
// 		Player: c.Location,
// 	}
// 	for client := range s.clients {
// 		if client != c {
// 			message := Envelope{
// 				Client:      strconv.Itoa(c.ID),
// 				MessageCode: EEvent,
// 				Message:     e,
// 			}
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				delete(s.clients, client)
// 			}
// 		}
// 	}
// }

// func (s *Server) LeavedPlayer(c *Client) {
// 	e := Event{
// 		Code:   "12",
// 		Id:     strconv.Itoa(c.ID),
// 		Player: c.Location,
// 	}
// 	for client := range s.clients {
// 		if client != c {
// 			message := Envelope{
// 				Client:      strconv.Itoa(c.ID),
// 				MessageCode: EEvent,
// 				Message:     e,
// 			}
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				delete(s.clients, client)
// 			}
// 		}
// 	}
// }
