package client

import "../protocol"

type messageHandler func(string)

type ChatClient interface {
	Dial(address string) error
	Start()
	Close()
	Send(command interface{}) error
	SendMessage(message string, receiver string) error
	SetName(message string) error
	Incoming() chan protocol.MessageCommand
	Errors() chan protocol.ErrorCommand
	Done() <-chan struct{}
}
