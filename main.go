package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

const (
	IP   = "127.0.0.1"
	PORT = 8080
)

type Room struct {
	Users []User
}

type User struct {
	Name       string
	ID         string
	Buffer     []byte
	Connection net.Conn
}

type Message struct {
	Sender  User
	Content string
}

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", IP, PORT))
	if err != nil {
		log.Fatal(err)
	}

	system := User{
		Name: "System",
	}

	room := Room{[]User{system}}
	newUser := make(chan User)
	newMessage := make(chan Message, 100)

	go listentonewConnections(listener, newUser)
	go readMessages(&room, newMessage)

	for {
		select {
		case msg := <-newMessage:
			sendMessageToRoom(room, msg)
			msg.Sender.Buffer = []byte("")
		case user := <-newUser:
			go appendConnectionToRoom(user, &room)
		}
	}
}

func appendConnectionToRoom(user User, room *Room) {
	greet(room.Users[0], &user)

	room.Users = append(room.Users, user)
}

func listentonewConnections(listener net.Listener, user chan User) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error while listining to connections %s\n", err)
		}

		user <- User{Connection: conn}
	}
}

func greet(from User, to *User) {
	name := make([]byte, 20)
	id := uuid.New()
	sendMessage(User{}, *to, []byte(fmt.Sprintf("Name: ")))
	n, err := to.Connection.Read(name)
	if err != nil {
		log.Println("Unable to get the name ", err)
	}
	to.Name = string(name[:n-1])
	to.ID = id.String()
	to.Buffer = make([]byte, 100)
	sendMessage(from, *to, []byte(fmt.Sprintf("Hi %s\n", to.Name)))
}

func readMessages(room *Room, done chan Message) {
	for {
		for i := 1; i < len(room.Users); i++ {
			connection := room.Users[i].Connection
			buffer := room.Users[i].Buffer

			n, err := connection.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Println(err)
			}
			done <- Message{
				Sender:  room.Users[i],
				Content: string(buffer[:n]),
			}
		}
	}
}

func sendMessageToRoom(room Room, message Message) {
	for i := 1; i < len(room.Users); i++ {
		currUser := room.Users[i]
		if message.Sender.ID != currUser.ID {
			sendMessage(message.Sender, room.Users[i], message.Sender.Buffer)
		}
	}
}

func sendMessage(from, to User, message []byte) {
	var msg = message
	if from.Name != "" {
		msg = []byte(fmt.Sprintf("%s: %s", from.Name, string(message)))
	}
	to.Connection.Write(msg)
}
