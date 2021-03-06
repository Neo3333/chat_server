package protocol

/**
version 1.0
*/

import (
	"bufio"
	"io"
	"log"
)

type CommandReader struct{
	reader *bufio.Reader
}

func NewCommandReader(reader io.Reader) *CommandReader{
	return &CommandReader{reader: bufio.NewReader(reader)}
}

func (r *CommandReader) Read() (interface{},error){
	commandName,err := r.reader.ReadString(' ')
	if err != nil{
		return nil,err
	}
	//log.Printf("commandName: %s, length %d",commandName,len(commandName))

	switch commandName {
	case "MESSAGE ":
		user,err := r.reader.ReadString(' ')
		if err != nil{
			return nil,err
		}
		time,err := r.reader.ReadString('*')
		if err != nil{
			return nil,err
		}
		message,err := r.reader.ReadString('\n')
		if err != nil{
			return nil,err
		}
		return MessageCommand{Name: user[:len(user)-1],Message:
			    				message[:len(message)-1],Time: time[:len(time)-1]},nil
	case "SEND ":
		to,err := r.reader.ReadString(' ')
		if err != nil{
			return nil,err
		}
		message,err := r.reader.ReadString('\n')
		if err != nil{
			return nil,err
		}
		return SendCommand{To: to[:len(to)-1], Message: message[:len(message)-1]},nil

	case "NAME ":
		name,err := r.reader.ReadString('\n')
		if err != nil{
			return nil,err
		}
		return NameCommand{Message: name[:len(name)-1]},err

	case "ERROR ":
		time,err := r.reader.ReadString('*')
		if err != nil{
			return nil,err
		}
		message,err := r.reader.ReadString('\n')
		if err != nil{
			return nil,err
		}
		return ErrorCommand{Message: message[:len(message)-1],Time: time[:len(time)-1]},nil

	default:
		log.Printf("Unknown command: %v", commandName)
	}
	return nil,UnknownCommand
}

func (r *CommandReader) ReadAll() ([]interface{},error){
	commands := []interface{}{}
	for {
		command, err := r.Read()
		if command != nil{
			commands = append(commands,command)
		}
		if err == io.EOF{
			break
		}else if err != nil{
			return commands,err
		}
	}
	return commands,nil
}
