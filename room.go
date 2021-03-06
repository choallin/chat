package main

import (
	"chat/trace"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type room struct {
	forward chan *message
	join    chan *client
	leave   chan *client
	clients map[*client]bool

	tracer trace.Tracer
}

func newRoom() *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),

		tracer: trace.Off(),
	}
}

func (r *room) allUsers() (users []string) {
	for client := range r.clients {
		users = append(users, client.userData["name"].(string))
	}
	return
}

func (r *room) run() {
	clientCode := regexp.MustCompile(`^ClientCode(\d){3}: (\w+) `)
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-r.join:
			msg := new(message)
			msg.Name = "server.socket"
			msg.Message = "Code002 " + client.userData["name"].(string)
			msg.Time = time.Now().Local().Format(time.Kitchen)
			fmt.Println("Neuer client: %v", client)
			fmt.Println("Nachricht: %v", msg)
			r.clients[client] = true
			for client := range r.clients {
				client.send <- msg
			}
		case client := <-r.leave:
			msg := NewMessageUserData(client.userData)
			delete(r.clients, client)
			close(client.send)
			for client := range r.clients {
				client.send <- msg
			}
		case msg := <-r.forward:
			fmt.Println("Message: %v", msg)
			for client := range r.clients {
				matches := clientCode.FindStringSubmatch(msg.Message)
				if len(matches) > 0 {
					group := matches[1]
					fmt.Println("RegExp Match: %v", matches)
					switch group {
					case "1":
						if client.userData["name"] != matches[2] {
							fmt.Println("Nur %s ist erlaubt", matches[2])
							continue
						}
					}
				}
				fmt.Println("Client: %v", client.userData["name"])
				select {
				case client.send <- msg:

				default:
					delete(r.clients, client)
					close(client.send)
					fmt.Println("Faild to deliver Message %v. Client closed", msg)
				}
			}
		case <-ticker.C:
			msg := new(message)
			for client := range r.clients {
				client.WriteMessage(msg)
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: messageBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("Can't get auth cookie:", err)
		return
	}
	cookieMap := map[string]interface{}{"name": authCookie.Value}
	client := &client{
		socket:   socket,
		send:     make(chan *message, messageBufferSize),
		room:     r,
		userData: cookieMap,
	}
	r.join <- client
	msg := new(message)
	msg.Name = "server.socket"
	allUsers := r.allUsers()

	msg.Message = "Code003 " + strings.Join(allUsers, ";")
	msg.Time = time.Now().Local().Format(time.Kitchen)
	client.send <- msg
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
