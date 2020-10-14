package server

type Command interface {
	Execute(*Conn, string)
}

type commandMap map[string]Command

var (
	commands = commandMap{
		"LOGIN": commandLogin{},
		"SEND": commandSend{},
		"QUIT": commandQuit{},
	}
)

type commandLogin struct {}

func (cmd commandLogin) Execute(conn *Conn, param string)  {
	
}

type commandSend struct {}

func (cmd commandSend) Execute(conn *Conn, param string)  {

}

type commandQuit struct {}

func (cmd commandQuit) Execute(conn *Conn, param string)  {

}