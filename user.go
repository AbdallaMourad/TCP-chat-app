package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/google/uuid"
)

type User struct {
	Name       string
	ID         string
	Buffer     []byte
	Connection net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	user := &User{
		Name:       name,
		Connection: conn,
	}

	user.GenerateNewID()
	user.AllocateBuffer(100)

	return user
}

func (u *User) AllocateBuffer(size uint32) {
	u.Buffer = make([]byte, size)
}

func (u *User) GenerateNewID() {
	u.ID = uuid.NewString()
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) CreateNewReader(room *Room) {
	go func() {
		for {
			n, err := u.Connection.Read(u.Buffer)
			if err != nil {
				if err == io.EOF {
					room.RemoveUser(u)
					room.SendMessageToRoom(Message{
						Sender:  *room.Admin,
						Content: fmt.Sprintf("%s has left the room\n", u.Name),
					})
					return
				}
				log.Println(err)
			}

			room.SendMessageToRoom(Message{
				Sender:  *u,
				Content: string(u.Buffer[:n]),
			})
		}
	}()
}
