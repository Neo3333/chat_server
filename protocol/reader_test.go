package protocol_test

import (
	"reflect"
	"strings"
	"testing"
	"../protocol"
)

func TestReaderCommand(t *testing.T){
	tests := []struct{
		input   string
		results []interface{}
	}{
		{
			"SEND test\n",
			[]interface{}{
				protocol.SendCommand{Message: "test"},
			},
		},
		{
			"MESSAGE user1 2020-02-02 00:00:00*hello\nMESSAGE user2 2006-01-02 15:04:05*world\n",
			[]interface{}{
				protocol.MessageCommand{Name: "user1",Message: "hello", Time: "2020-02-02 00:00:00"},
				protocol.MessageCommand{Name: "user2",Message: "world", Time: "2006-01-02 15:04:05"},
			},
		},
	}

	for _,test := range tests{
		reader := protocol.NewCommandReader(strings.NewReader(test.input))
		results, err := reader.ReadAll()

		t.Log(results)

		if err != nil{
			t.Errorf("Unable to read command, error %v", err)
		}else if !reflect.DeepEqual(results,test.results){
			t.Errorf("Command output is not the same: %v %v", results, test.results)
		}
	}
}
