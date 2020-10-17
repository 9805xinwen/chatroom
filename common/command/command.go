package commands

import "flag"

type Command interface {
	GetFlagSet() flag.FlagSet

	Execute(str string) error
}
