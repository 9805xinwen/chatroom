package commands

import (
	"chatroom/common/utils"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

var EchoCommand Command = CreateDefaultCommand("echo", reflect.TypeOf(EchoData{}), EchoRun)

type EchoData struct {
	Times   int64  `name:"t"       value:"1"   usage:"输出的重复次数"`
	Content string `name:"content" value:""    usage:"输出的内容"`
}

// echo -content "hello!" -t 100
func EchoRun(params Params) error {
	data := params.Info.(*EchoData)

	if data.Content == "" {
		if len(params.Args) > 0 {
			data.Content = params.Args[0]
		}
	}

	outputStr, _ := utils.DoubleQuotedStringsMarch(data.Content)
	outputStr = strings.ReplaceAll(outputStr, "%20", " ")

	var i int64 = 0
	for i = 0; i < data.Times; i = i + 1 {
		fmt.Println(outputStr)
	}

	return nil
}

//////////////////////////////////////////////////////////////
//                    echo 命令使用测试                       //
//////////////////////////////////////////////////////////////

func TestSayHelloCommand(t *testing.T) {
	cmd := "echo -content \"hello!\" "
	bundle := map[string]interface{}{
		"Connect": "connecting", //net.Conn
	}
	EchoCommand.Execute(cmd, bundle)

}
