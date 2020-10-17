package commands

import (
	"fmt"
	"reflect"
	"testing"
)

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

var EchoCommand Command = CreateDefaultCommand("echo", reflect.TypeOf(SayData{}), SayRun)

type SayData struct {
	Times   int64  `name:"t"       value:"1"   usage:"输出的重复次数"`
	Content string `name:"content" value:""    usage:"输出的内容"`
}

func SayRun(params Params) error {
	data := params.Info.(*SayData)

	if data.Content == "" {
		if len(params.Args) > 0 {
			data.Content = params.Args[0]
		}
	}

	var i int64 = 0
	for i = 0; i < data.Times; i = i + 1 {
		fmt.Println(data.Content)
	}

	return nil
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////
func TestSayHelloCommand(t *testing.T) {
	cmd := "echo \"你在干什么？\""
	EchoCommand.Execute(cmd, nil)
}
