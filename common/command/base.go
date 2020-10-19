package commands

import (
	"flag"
	"strings"
)

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type Runner interface {
	Run(params Params) error
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////

type Params struct {
	Str string

	Info interface{}
}

type ParamValue interface {
	flag.Getter
}

type ModelProvider interface {
	GetParamsModel() interface{}
}

//////////////////////////////////////////////////////////////
//                                                          //
//////////////////////////////////////////////////////////////
type CommandBase struct {
	Command

	runner Runner

	flags flag.FlagSet

	modelProvider ModelProvider
}

func NewCommand(runner Runner, flags flag.FlagSet) Command {
	return &CommandBase{runner: runner, flags: flags}
}

func (cmd *CommandBase) Execute(str string) error {

	//构建模型
	model := cmd.modelProvider.GetParamsModel()

	//复制flags

	//反射绑定

	//对命令进行解析
	err := cmd.flags.Parse(strings.Fields(str))
	if err != nil {
		return err
	}

	//设置参数
	params := Params{Str: str, Info: model}

	//将参数传入runner
	return cmd.runner.Run(params)
}
