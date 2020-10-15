package server

import "flag"

type Command interface {
	Execute(*Conn, string)
}

type commandMap map[string]Command

var (
	commands = commandMap{
		LoginCommandName: commandLogin{},
		"SEND":           commandSend{},
		"QUIT":           commandQuit{},
	}
)

type commandLogin struct{}

const LoginCommandName string = "LOGIN"

var loginFlag *flag.FlagSet = flag.NewFlagSet(LoginCommandName, flag.ContinueOnError)

func (cmd commandLogin) Execute(conn *Conn, param string) {

}

type commandSend struct{}

func (cmd commandSend) Execute(conn *Conn, param string) {

}

type commandQuit struct{}

func (cmd commandQuit) Execute(conn *Conn, param string) {

}
