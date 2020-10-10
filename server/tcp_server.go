package server

import (
	"../protocol"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

type client struct{
	conn        net.Conn
	name        string
	writer      *protocol.CommandWriter
}

type TcpChatServer struct {
	listener    net.Listener
	clients     map[string]*client
	lock        *sync.Mutex
}

var(
	UnknownClient = errors.New("Unknown client")
	DuplicateName = errors.New("Duplicate Name")
)

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
	for _, client := range s.clients {
		if client.name == name {
			return client.writer.Write(command)
		}
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
	s.clients[client.name] = client
	return client
}

func (s *TcpChatServer) remove(client *client){
	s.lock.Lock()
	defer s.lock.Unlock()
	c := s.clients[client.name]
	if (c == nil){
		//TODO 增加错误处理机制
		//panic("Data Corruption")
		log.Fatal("Data Corruption")
		return;
	}

	delete(s.clients,client.name)
	log.Printf("Closing connection from %v",client.conn.RemoteAddr().String())
	_ = client.conn.Close()
}

func (s *TcpChatServer) serve(client *client){
	cmdReader := protocol.NewCommandReader(client.conn)
	defer s.remove(client)
	for{
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF{
			log.Printf("Unknown command from: %s with error %v",
				client.conn.RemoteAddr().String(),err)
			_ = client.writer.Write(protocol.ErrorCommand{
				Message: err.Error(),
			})
		}
		if cmd != nil{
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name: client.name,
				})
			case protocol.NameCommand:
				s.changeName(client, v.Message)
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
	c := s.clients[newName]
	if c != nil{
		newName += "@" + client.conn.RemoteAddr().String()
	}
	delete(s.clients,client.name)
	s.clients[newName] = client
	client.name = newName
}

