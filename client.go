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
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.Time = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			// PingMessage == 9
			if err := c.socket.WriteMessage(9, []byte{}); err != nil {
				fmt.Println("PingMessage not availiable")
			}
		default:
			for msg := range c.send {
				if err := c.socket.WriteJSON(msg); err != nil {
					fmt.Println("Error: %v", err)
					break
				}
			}
			c.socket.Close()
		}
	}
}
