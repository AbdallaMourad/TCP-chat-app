package main

import (
	"fmt"
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
	userJoiningChannel := make(chan User)

	go listentonewConnections(listener, userJoiningChannel)

	for {
		select {
		case user := <-userJoiningChannel:
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
