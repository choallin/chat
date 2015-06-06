package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

const (
	pongWait       = 60 * time.Second
	writeWait      = 10 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	PingMessage    = 9
)

func (c *client) WriteMessage(msg *message) {
	if err := c.socket.WriteMessage(PingMessage, []byte(msg.Message)); err != nil {
		fmt.Println("PingMessage not availiable")
		fmt.Println("Error: %v", err)
	} else {
		fmt.Println("Ping mesage sent")
	}
}

func (c *client) read() {
	c.socket.SetReadLimit(maxMessageSize)
	c.socket.SetReadDeadline(time.Now().Add(pongWait))
	c.socket.SetPongHandler(func(string) error { c.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.Time = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			c.user_lost()
			fmt.Println("Error in read Method")
			fmt.Println("Error: %v", err)
			break
		}
	}
	c.socket.Close()
}

func (c *client) user_lost() {
	msg := NewMessageUserData(c.userData)
	fmt.Println("Send user lost message. %v", msg)
	c.socket.WriteJSON(msg)
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
