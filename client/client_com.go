package client

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)


type (
	Command interface {
		RequireParam() bool
		CommandFormat(str string,reqpara bool)
		Execute()
	}
)

type commandMap map[string]Command

var (
	Commands = commandMap{
		"login": commandlogin{},
		"send":  commandsend{},
		"quit":  commandquit{},
	})

type commandlogin struct{}

func (com commandlogin) RequireParam() bool{
	return true
}

func (com commandlogin) CommandFormat(str string,reqpara bool) {
	if reqpara{
		params := strings.SplitN(strings.Trim(str, "\r\n"), " ", 2)
		len := len(params)
		if len == 1 {
			for _, r := range str {
				if !(unicode.IsDigit(r)) {
					fmt.Println("输入不合法")
					fmt.Printf("请输入(命令 参数):\n")
				}
			}
		}
		fmt.Println("输入合法")
		fmt.Printf("请输入(命令 参数):\n")
	}
}

func (com commandlogin) Execute() {

}

type commandsend struct{}

func (com commandsend) RequireParam() bool{
	return true
}

func(com commandsend) CommandFormat(str string,reqpara bool) {
	if reqpara{
		params := strings.SplitN(strings.Trim(str, "\r\n"), " ", 2)
		len := len(params)
		if len == 2{
			username := params[0]
			msg := params[1]
			if username != "" && msg != ""{
				fmt.Println("输入合法")
				fmt.Printf("请输入(命令 参数):\n")
			}
		}else{
			fmt.Println("输入不合法")
			fmt.Printf("请输入(命令 参数):\n")
		}
	}
}
func (com commandsend) Execute() {}

type commandquit struct{}

func (com commandquit) RequireParam() bool{
	return false
}

func (com commandquit) CommandFormat(param string,reqpara bool) {
	if !reqpara{
		if param == ""{
			fmt.Println("退出成功！")
			com.Execute()
		}
	}
}

func (com commandquit) Execute() {
	os.Exit(-1)
}

