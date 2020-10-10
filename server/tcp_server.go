package server

import (
	"io"
	"log"
	"net"
	"../protocol"
	"sync"
)

type client struct{
	conn        net.Conn
	name        string
	writer      *protocol.CommandWriter
}

type TcpChatServer struct {
	listener    net.Listener
	clients     []*client
	lock        *sync.Mutex
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
		if err != nil{
			log.Printf("Read error %v",err)
		}
	}
	return nil
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
	s.clients = append(s.clients)
	return client
}

func (s *TcpChatServer) remove(client *client){
	s.lock.Lock()
	defer s.lock.Unlock()
	for i,check := range s.clients{
		if check == client{
			s.clients = append(s.clients[:i],s.clients[i+1:]...)
		}
	}
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
				client.name = v.Message
			}
		}
		if err == io.EOF{
			break
		}
	}
}

