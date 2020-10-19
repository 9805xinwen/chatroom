package commands

type Command interface {
	Execute(str string, bundle map[string]interface{}) error
}
