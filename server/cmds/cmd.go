package cmds

import "chatroom/common/commands"

//命令行处理集
var CommandMap map[string]commands.Command = map[string]commands.Command{
	LoginCommandName: LoginCommand,
	SendCommandName:  SendCommand,
	QuitCommandName:  QuitCommand,
}
