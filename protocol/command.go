package protocol

/**
version 1.0
 */

import (
	"errors"
)

var(
	UnknownCommand = errors.New("Unknown Command")
)

type SendCommand struct{
	Message string
	To      string
}

type NameCommand struct{
	Message string
}

type MessageCommand struct{
	Name    string
	Message string
	Time    string
}

type ErrorCommand struct{
	Message string
	Time    string
}




