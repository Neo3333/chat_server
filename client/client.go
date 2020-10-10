package client

import "../protocol"

type messageHandler func(string)

type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SendMessage(message string) error
	SetName(message string) error
	Incoming() chan protocol.MessageCommand
}
