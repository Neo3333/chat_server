package client

import (
	"../protocol"
	"io"
	"log"
	"net"
	"time"
)

type TcpChatClient struct {
	conn net.Conn
	name string
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	incoming chan protocol.MessageCommand
	errors   chan protocol.ErrorCommand
}

func NewClient() *TcpChatClient{
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
		errors: make(chan protocol.ErrorCommand),
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

func (c *TcpChatClient) Start(){
	for{
		cmd,err := c.cmdReader.Read()
		if err == io.EOF{
			break
		}else if err != nil{
			log.Printf("Read error %v", err)
		}

		if cmd != nil{
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			case protocol.ErrorCommand:
				//log.Printf("Error from the server %s",v.Message)
				c.errors <- v
			default:
				c.errors <- protocol.ErrorCommand{
				Message: "Unknown Command",
				Time: time.Now().Format("2006-01-02 15:04:05"),
				}
			}
		}
	}
}

func (c *TcpChatClient) Close() {
	_ = c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *TcpChatClient) Errors() chan protocol.ErrorCommand  {
	return c.errors
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{Message: name})
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
		To: "*",
	})
}

func (c *TcpChatClient) SendMessagePrivate(message string, receiver string) error{
	return c.Send(protocol.SendCommand{
		Message: message,
		To: receiver,
	})
}
