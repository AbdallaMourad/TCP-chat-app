package main

import (
	"fmt"
	"log"
)

type Room struct {
	Users []User
}

func NewRoom(system User) *Room {
	return &Room{
		Users: []User{system},
	}
}

func (r *Room) JoinRoom(newUser User) {
	system := r.Users[0]

	getNameMessage := NewMessage(system, "Name: ")

	name := r.SendSystemMessage(getNameMessage, newUser)
	newUser.SetName(name)

	newUserJoinedRoomMessage := NewMessage(system, fmt.Sprintf("%s has joined the room\n", name))

	r.Broadcast(newUserJoinedRoomMessage)

	r.Users = append(r.Users, newUser)
}

func (r *Room) SendSystemMessage(message Message, user User) string {
	response := make([]byte, 20)
	user.Connection.Write([]byte(message.Content))

	n, err := user.Connection.Read(response)
	if err != nil {
		log.Println("Unable to get the name ", err)
	}

	return string(response[:n-1])
}

func (r *Room) Broadcast(message Message) {
	r.SendMessageToRoom(message)
}

func (r *Room) SendMessageToRoom(message Message) {
	for i := 1; i < len(r.Users); i++ {
		currUser := r.Users[i]
		if message.Sender.ID != currUser.ID {
			currUser.Connection.Write(message.GetFormattedMessage())
		}
	}
}
