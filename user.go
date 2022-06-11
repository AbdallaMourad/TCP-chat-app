package main

import (
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
