package main

import (
	"../../client"
	"flag"
	"log"
	"../../tui"
)

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
	tui.StartUi(c)
}

