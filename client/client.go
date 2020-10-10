package client

import "../protocol"

type messageHandler func(string)

type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SetName(message string) error
	Incomming() chan protocol.MessageCommand
}
