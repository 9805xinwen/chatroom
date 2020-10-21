package cmds

import (
	"chatroom/common/commands"
	"chatroom/server/services/message"
	"chatroom/server/services/users"
)

//命令行处理集
var CommandMap map[string]commands.Command = map[string]commands.Command{
	//    命令名        |     命令实例     |    备注
	LoginCommandName    :  LoginCommand    ,  // "login"
	SendCommandName     :  SendCommand     ,  // "send"
	QuitCommandName     :  QuitCommand     ,  // "quit"
	RegisterCommandName :  RegisterCommand ,  // "register"
}

var GlobalOnlineService users.Online
var GlobalUserService users.Users
var GlobalMassageService message.MassageService