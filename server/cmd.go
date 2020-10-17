package server

type Command interface {
	Execute(*Conn, string)
	RequireAuth() bool
}

type commandMap map[string]Command

var (
	commands = commandMap{
		"LOGIN": commandLogin{},
		"SEND":  commandSend{},
		"QUIT":  commandQuit{},
	}
)

type commandLogin struct{}

func (cmd commandLogin) RequireAuth() bool {
	return false
}

func (cmd commandLogin) Execute(conn *Conn, param string) {
	conn.reqUser = param
	if name, ok := conn.server.Auth.CheckId(conn.reqUser); ok {
		conn.user = name
		conn.writeMessage("ok")
		return
	}
	return
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
