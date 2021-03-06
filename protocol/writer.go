package protocol

import (
	"fmt"
	"io"
)

/**
version 1.0
*/

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter{
	return &CommandWriter{writer: writer}
}

func (w *CommandWriter) writeString(msg string) error{
	_, err := w.writer.Write([]byte(msg))
	return err
}

func (w *CommandWriter) Write(command interface{}) error{
	var err error
	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v %v\n", v.To, v.Message))
	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Message))
	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v %v*%v\n", v.Name, v.Time, v.Message))
	case ErrorCommand:
		err = w.writeString(fmt.Sprintf("ERROR %v*%v\n",v.Time, v.Message))
	default:
		err = UnknownCommand
	}
	return err
}