package main

import (
	"../../client"
	"flag"
	"log"
	"../../tui"
)

/**
version 1.0
*/

func main(){
	address := flag.String("server","","server address to connect")
	flag.Parse()

	c := client.NewClient()
	err := c.Dial(*address)

	if err != nil{
		log.Fatal(err)
	}
	defer c.Close()

	go c.Start()

	go func() {
		arr := make([]byte, len(*address))
		copy(arr, *address)
		ip := string(arr)
		tui.Server_ip <- ip[:len(ip)-5]
	}()

	tui.StartUi(c)
}

