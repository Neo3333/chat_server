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
	//clients     []*client
	clients     map[*client]string
	names       map[string]bool
	lock        *sync.Mutex
}

var(
	UnknownClient = errors.New("Unknown client")
	DuplicateName = errors.New("Duplicate Name")
)

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		lock: &sync.Mutex{},
		clients: make(map[*client]string),
		names: make(map[string]bool),
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
	for client,_ := range s.clients{
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
	for client, _ := range s.clients {
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
	s.clients[client] = client.name
	s.names[client.name] = true
	return client
}

func (s *TcpChatServer) remove(client *client){
	s.lock.Lock()
	defer s.lock.Unlock()
	_,ok := s.clients[client]
	if (!ok){
		//TODO 增加错误处理机制
		//panic("Data Corruption")
		log.Fatal("Data Corruption")
		return;
	}

	delete(s.clients,client)
	delete(s.names,client.name)
	log.Printf("Closing connection from %v",client.conn.RemoteAddr().String())
	_ = client.conn.Close()
}

func (s *TcpChatServer) serve(client *client){
	cmdReader := protocol.NewCommandReader(client.conn)
	defer s.remove(client)
	for{
		cmd, err := cmdReader.Read()
		if err != nil && err != io.EOF{
			log.Printf("Read error: %v",err)
		}
		if cmd != nil{
			switch v := cmd.(type) {
			case protocol.SendCommand:
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name: client.name,
				})
			case protocol.NameCommand:
				err := s.changeName(client, v.Message)
				if err != nil{
					client.writer.Write(protocol.ErrorCommand{
						Message: err.Error(),
					})
				}
			}
		}
		if err == io.EOF{
			break
		}
	}
}

func (s *TcpChatServer) changeName(client *client, newName string) error{
	s.lock.Lock()
	defer s.lock.Unlock()
	_,ok := s.names[newName]
	if ok{
		return DuplicateName
	}
	delete(s.names,client.name)
	s.names[newName] = true

	client.name = newName
	s.clients[client] = newName
	return nil
}

