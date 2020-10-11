package server

import (
	"../protocol"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
	"../utils"
)
/**
version 1.0
*/

type client struct{
	conn        net.Conn
	name        string
	writer      *protocol.CommandWriter
}

type TcpChatServer struct {
	listener    net.Listener
	clients     map[string]*client
	lock        *sync.Mutex
	generator   *utils.UUIDGenerator
}

var UnknownClient = errors.New("Unknown client")
var loadGeneratorOnce sync.Once
const SYSTEM = "system"

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		lock: &sync.Mutex{},
		clients: make(map[string]*client),
	}
}

func (s *TcpChatServer) Listen(address string) error{
	l,err := net.Listen("tcp",address)
	if err == nil{
		s.listener = l
	}
	log.Printf("Listening on %v",address)
	return err
}

func (s *TcpChatServer) Close(){
	_ = s.listener.Close()
	s.generator.Close()
}

func (s *TcpChatServer) Start(){
	for{
		conn, err := s.listener.Accept()
		if err != nil{
			log.Print(err)
		}else{
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *TcpChatServer) Broadcast(command interface{}) error{
	for _,client := range s.clients{
		err := client.writer.Write(command)
		//TODO 更好的错误处理机制
		if err != nil{
			log.Printf("Connection lost on %s with error %v",
				client.conn.RemoteAddr(),err)
		}
	}
	return nil
}

func (s *TcpChatServer) Send(name string, command interface{}) error {
	if c := s.clients[name]; c != nil{
		return c.writer.Write(command)
	}
	return UnknownClient
}

func (s *TcpChatServer) accept(conn net.Conn) *client{
	log.Printf("Accepting connection from %v, total clients: %v",
		conn.RemoteAddr().String(), len(s.clients)+1)
	s.lock.Lock()
	defer s.lock.Unlock()
	client := &client{
		conn: conn,
		name: conn.RemoteAddr().String(),
		writer: protocol.NewCommandWriter(conn),
	}
	go s.Broadcast(protocol.MessageCommand{
		Message: fmt.Sprintf("User@%s arrived",conn.RemoteAddr().String()),
		Name: SYSTEM,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	})
	s.clients[client.name] = client
	return client
}

func (s *TcpChatServer) remove(client *client){
	s.lock.Lock()
	defer s.lock.Unlock()
	c := s.clients[client.name]
	if (c == nil){
		//TODO 增加错误处理机制
		log.Fatal("Data Corruption")
		return;
	}
	go s.Broadcast(protocol.MessageCommand{
		Message: fmt.Sprintf("User<%s> left Skynet",client.name),
		Name: SYSTEM,
		Time: time.Now().Format("2006-01-02 15:04:05"),
	})
	delete(s.clients,client.name)
	log.Printf("Closing connection from %v",client.conn.RemoteAddr().String())
	_ = client.conn.Close()
}

func (s *TcpChatServer) serve(cli *client){
	cmdReader := protocol.NewCommandReader(cli.conn)
	defer s.remove(cli)
	for{
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF{
			log.Printf("Unknown command from: %v with error %v",
				cli.conn.RemoteAddr().String(),err)
			_ = cli.writer.Write(protocol.ErrorCommand{
				Message: err.Error(),
			})
		}
		if cmd != nil{
			switch v := cmd.(type) {
			case protocol.SendCommand:
				msg := protocol.MessageCommand{
					Message: v.Message,
					Name: cli.name,
					Time: time.Now().Format("2006-01-02 15:04:05"),
				}
				switch v.To {
				case "*":
					go s.Broadcast(msg)
				default:
					msg.Name = "[PRIVATE]"+msg.Name
					go func(c *client, rec string, msg interface{}) {
						err := s.Send(v.To, msg)
						if err != nil{
							_ = c.writer.Write(protocol.ErrorCommand{
								Message: err.Error(),
								Time:    time.Now().Format("2006-01-02 15:04:05"),
							})
						}else{
							_ = c.writer.Write(msg)
						}
					}(cli,v.To,msg)
				}
			case protocol.NameCommand:
				s.changeName(cli, v.Message)
			}
		}
		if err == io.EOF{
			break
		}
	}
}

func (s *TcpChatServer) changeName(client *client, newName string){
	s.lock.Lock()
	defer s.lock.Unlock()
	name := newName
	c := s.clients[name]
	for name=="system" || c != nil{
		loadGeneratorOnce.Do(s.startGenerator)
		name = newName + "#" + s.generator.Get()
		c = s.clients[name]
	}
	delete(s.clients,client.name)
	s.clients[name] = client
	client.name = name
}

func (s *TcpChatServer) startGenerator(){
	s.generator = utils.NewUUIDGenerator("")
}


