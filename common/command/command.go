package commands

var commandList = []string{"LOGIN", "SEND", "QUIT"}

type Command struct {

	Blocker
	Executer
}