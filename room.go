package main

import (
	"fmt"
	"log"
)

type Room struct {
	Admin *User
	Users map[*User]bool
}

func NewRoom(admin *User) *Room {
	return &Room{
		Admin: admin,
		Users: make(map[*User]bool),
	}
}

func (r *Room) JoinRoom(newUser *User) {
	getNameMessage := NewMessage(*r.Admin, "Name: ")

	name := r.SendSystemMessage(getNameMessage, newUser)
	newUser.SetName(name)

	newUserJoinedRoomMessage := NewMessage(*r.Admin, fmt.Sprintf("%s has joined the room\n", name))

	r.Broadcast(newUserJoinedRoomMessage)

	r.Users[newUser] = true

	newUser.CreateNewReader(r)
}

func (r *Room) SendSystemMessage(message Message, user *User) string {
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
	for user := range r.Users {
		if message.Sender.ID != user.ID {
			user.Connection.Write(message.GetFormattedMessage())
		}
	}
}

func (r *Room) RemoveUser(user *User) {
	if _, ok := r.Users[user]; ok {
		delete(r.Users, user)
	}
}
