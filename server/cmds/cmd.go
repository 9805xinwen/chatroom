package cmds

import (
	"chatroom/common/commands"
	"chatroom/server/services/message"
	"chatroom/server/services/users"
)

//命令行处理集
var CommandMap map[string]commands.Command = map[string]commands.Command{
	LoginCommandName: LoginCommand,
	SendCommandName:  SendCommand,
	QuitCommandName:  QuitCommand,
}

var GlobalOnlineService users.Online
var GlobalUserService users.Users
var GlobalMassageService message.MassageService