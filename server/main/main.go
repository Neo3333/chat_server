package main

import (
"../../server"
	"log"
	"os"
)

/**
version 1.0
*/

func main() {
	var s server.ChatServer
	s = server.NewServer()
	err := s.Listen("127.0.0.1:8000")

	if err != nil{
		log.Print(err)
		os.Exit(30)
	}

	// start the server
	s.Start()
}