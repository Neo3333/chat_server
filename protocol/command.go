package protocol

import "errors"

var(
	UnknownCommand = errors.New("Unknown Command")
)

type SendCommand struct{
	Message string
}

type NameCommand struct{
	Message string
}

type MessageCommand struct{
	Name    string
	Message string
}

type ErrorCommand struct{
	Message string
}



