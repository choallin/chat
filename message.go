package main

import "time"

type message struct {
	Message string
	Name    string
	Time    string
}

func NewMessageUserData(userData map[string]interface{}) *message {
	msg := new(message)
	if name, ok := userData["name"].(string); ok {
		msg.Message = "Code001 " + name
		msg.Name = "server.socket"
		msg.Time = time.Now().Local().Format(time.Kitchen)
		return msg
	}
	return nil
}
