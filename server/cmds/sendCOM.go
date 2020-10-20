package cmds

import (
	"chatroom/common/commands"
	"chatroom/common/utils"
	"chatroom/server/services/message"
	"errors"
	"net"
	"reflect"
)

////////////////////////////////////////////////////////////////////////
//                            Send 命令定义                           //
//--------------------------------------------------------------------//
// [命令名称] : send                                                  //
// [命令参数] :                                                       //
//             -msg                          发送的信息               //
//             -user                         发送给某个用户           //
//--------------------------------------------------------------------//
// 使用案例:                                                          //
// send -user USERNAME -msg "hello!"                                  //
////////////////////////////////////////////////////////////////////////

const SendCommandName string = "send"

var SendCommand commands.Command = commands.CreateDefaultCommand(SendCommandName, reflect.TypeOf(SendData{}), SendRun)

////////////////////////////////////////////////////////////////////////
//                        主要命令参数结构体定义                      //
//--------------------------------------------------------------------//
// 定义结构体:                                                        //
//        SendData                           登录数据结构体           //
//--------------------------------------------------------------------//
// SendData                                                           //
// [公开属性] :                                                       //
//   - Username                                字符串 | 用户名        //
//   - Massage                                 字符串 | 发送的信息    //
// [私有属性] : -无-                                                  //
// [构造函数] : -无-                                                  //
// [公开函数] : -无-                                                  //
// [私有函数] : -无-                                                  //
////////////////////////////////////////////////////////////////////////

type SendData struct {
	Username string `name:"user" value:"" usage:"发送给某个用户"`
	Massage  string `name:"msg"  value:"" usage:"发送的信息"`
}

////////////////////////////////////////////////////////////////////////
//                        主要函数(runner)实现                        //
//--------------------------------------------------------------------//
// 实现函数:                                                          //
//        SendRun(params commands.Params)        发送信息处理         //
//--------------------------------------------------------------------//
// 使用的内部的参数结构体(Params.Info属性对应的结构体) ： SendRun     //
////////////////////////////////////////////////////////////////////////

func SendRun(params commands.Params) error {
	//获取解析数据
	data := params.Info.(SendData)

	//检查是否存在发送目标用户
	_, err := GlobalUserService.GetId(data.Username)
	if err != nil {
		return errors.New("目标用户不存在！")
	}

	//获取发送者UserId
	userId := params.Bundle[UserId].(string)

	//从userId获取username
	username, err := GlobalUserService.GetName(userId)
	if err != nil {
		return err
	}

	//获取发送者的连接
	fromConn := params.Bundle[Connect].(*net.Conn)

	//获取发送连接,不在线就设置为nil
	var toConn *net.Conn
	if !GlobalOnlineService.OnlineCheckByUserName(data.Username) {
		toConn = nil
	} else {
		toConn = GlobalOnlineService.QueryConnByUserName(data.Username)
	}

	//发送数据内容处理
	content, err := utils.DoubleQuotedStringsMarch(data.Massage)
	if err != nil {
		return err
	}

	msg := message.Massage{
		FromUser: username,
		ToUser:   data.Username,
		FromConn: fromConn,
		ToConn:   toConn,
		Content:  content,
	}
	GlobalMassageService.Send(msg)

	return nil
}
