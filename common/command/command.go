package commands

import "flag"

var commandList = []string{"LOGIN", "SEND", "QUIT"}

type Command interface {
	GetFlagSet() flag.FlagSet

	Execute(str string) error
}
