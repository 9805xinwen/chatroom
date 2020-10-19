package server

type Command interface {
	//Execute(*net.Conn, string)
	//RequireAuth() bool
}

type commandMap map[string]Command

var (
	commands = commandMap{
		//"LOGIN": commandLogin{},
		//"REGISTER": commandRegister{},
		//"SEND":  commandSend{},
		//"QUIT":  commandQuit{},
	}
)

