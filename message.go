package main

import "fmt"

type Message struct {
	Sender  User
	Content string
}

func NewMessage(sender User, content string) Message {
	return Message{
		Sender:  sender,
		Content: content,
	}
}

func (m Message) GetFormattedMessage() []byte {
	return []byte(fmt.Sprintf("%s: %s", m.Sender.Name, m.Content))
}
