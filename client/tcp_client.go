package client

import (
	"net"
	"../protocol"
)

type TcpChatClient struct {
	conn net.Conn
	name string
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	incoming chan protocol.MessageCommand
}

func NewTcpChatClient() *TcpChatClient{
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
	}
}

func (c *TcpChatClient) Dial(address string) error{
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
	}
	c.cmdReader = protocol.NewCommandReader(conn)
	c.cmdWriter = protocol.NewCommandWriter(conn)
	return err
}
