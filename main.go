package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

const (
	IP   = "127.0.0.1"
	PORT = 8080
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", IP, PORT))
	if err != nil {
		log.Fatal(err)
	}

	system := NewUser("System", nil)

	room := NewRoom(*system)
	newUser := make(chan User)
	newMessage := make(chan Message, 100)

	go listentonewConnections(listener, newUser)
	go readMessages(room, newMessage)

	for {
		select {
		case msg := <-newMessage:
			go room.SendMessageToRoom(msg)
			msg.Sender.Buffer = []byte("")
		case user := <-newUser:
			go room.JoinRoom(user)
		}
	}
}

func listentonewConnections(listener net.Listener, user chan User) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error while listining to connections %s\n", err)
		}

		newUser := NewUser("", conn)
		user <- *newUser
	}
}

func readMessages(room *Room, done chan Message) {
	arr := make([]func(), len(room.Users))

	for {
		for i := 1; i < len(room.Users); i++ {
			if len(arr) < len(room.Users) {
				arr = append(arr, func() {
					go readMessage(room.Users[i], done)
				})
				arr[i]()
			}
		}
	}
}

func readMessage(user User, done chan Message) {
	n, err := user.Connection.Read(user.Buffer)
	if err != nil {
		if err == io.EOF {
			return
		}
		log.Println(err)
	}
	done <- Message{
		Sender:  user,
		Content: string(user.Buffer[:n]),
	}
}
