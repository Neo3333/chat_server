package protocol_test

import (
	"bytes"
	"log"
	"testing"
	"../protocol"
)

func TestWriteCommand(t *testing.T){
	tests := []struct{
		commands []interface{}
		result   string
	}{
		{
			[]interface{}{
				protocol.SendCommand{"Hello"},
			},
			"SEND Hello\n",
		},
		{
			[]interface{}{
				protocol.MessageCommand{Message: "world",
					Name: "michael",
					Time: "2006-01-02 15:04:05"},
			},
			"MESSAGE michael 2006-01-02 15:04:05*world\n",
		},
	}

	buf := new(bytes.Buffer)
	for _,test := range tests{
		buf.Reset()
		cmdWriter := protocol.NewCommandWriter(buf)

		for _,cmd := range test.commands{
			err := cmdWriter.Write(cmd)
			if (err != nil){
				t.Errorf("Unable to write command %v",cmd)
			}
		}

		if buf.String() != test.result{
			log.Println(len(buf.String()))
			log.Println(len(test.result))
			t.Errorf("Command output is not the same %v %v", buf.String(),test.result)
		}
	}
}
