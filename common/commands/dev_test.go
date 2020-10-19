package commands

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
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

func EchoRun(params Params) error {
	data := params.Info.(*EchoData)

	if data.Content == "" {
		if len(params.Args) > 0 {
			data.Content = params.Args[0]
		}
	}

	outputStr, _ := DoubleQuotedStringsMarch(data.Content)
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
	//echo "hello%20world%20!"
	cmd := "echo \"Hello%20world%20!\" "
	EchoCommand.Execute(cmd, nil)
}
