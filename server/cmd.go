package server

import (
	"fmt"
	"log"
	"regexp"
)

type Command interface {
	Execute(*Conn, string)
	RequireAuth() bool
}

type commandMap map[string]Command

var (
	commands = commandMap{
		"LOGIN": commandLogin{},
		"REGISTER": commandRegister{},
		"SEND":  commandSend{},
		"QUIT":  commandQuit{},
	}
)

type commandLogin struct{}

func (cmd commandLogin) RequireAuth() bool {
	return false
}

func (cmd commandLogin) Execute(conn *Conn, param string) {
	if conn.userId != "" && conn.userName != "" {
		conn.writeMessage("已登录，请勿重复login")
	}
	conn.reqUserId = param
	name, err := conn.server.User.GetName(conn.reqUserId)
	if err != nil {
		log.Print(err)
		conn.writeMessage("错误：id不存在")
		return
	}
	conn.userId = conn.reqUserId
	conn.userName = name
	conn.writeMessage("登陆成功")
}

type commandRegister struct{}

func (cmd commandRegister) RequireAuth() bool {
	return false
}

func (cmd commandRegister) Execute(conn *Conn, param string) {
	name := param
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]{4,16}$", name)
	if !matched {
		conn.writeMessage("400 用户名格式不正确")
		return
	}

	id, err := conn.server.User.Register(name)
	if err != nil {
		conn.writeMessage(err.Error())
		return
	}
	conn.userId = id
	conn.userName = name
	conn.writeMessage(fmt.Sprintf("注册成功，请牢记住您的唯一id: %s", id))
}

type commandSend struct{}

func (cmd commandSend) RequireAuth() bool {
	return true
}

func (cmd commandSend) Execute(conn *Conn, param string) {

}

type commandQuit struct{}

func (cmd commandQuit) RequireAuth() bool {
	return false
}

func (cmd commandQuit) Execute(conn *Conn, param string) {

}
