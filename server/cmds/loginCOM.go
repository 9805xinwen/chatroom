package cmds

import (
	"chatroom/common/commands"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
)

////////////////////////////////////////////////////////////////////////
//                           Login 命令定义                           //
//--------------------------------------------------------------------//
// [命令名称] : login                                                 //
// [命令参数] :                                                       //
//            -id                      [默认]用户id                   //
//--------------------------------------------------------------------//
// 使用案例:                                                          //
// login -id USERNAME                                                 //
// login USERNAME                                                     //
////////////////////////////////////////////////////////////////////////

const LoginCommandName string = "login"

var LoginCommand commands.Command = commands.CreateDefaultCommand(LoginCommandName, reflect.TypeOf(LoginData{}), LoginRun)

////////////////////////////////////////////////////////////////////////
//                        主要命令参数结构体定义                      //
//--------------------------------------------------------------------//
// 定义结构体:                                                        //
//        LoginData                          登录数据结构体           //
//--------------------------------------------------------------------//
// LoginData                                                          //
// [公开属性] :                                                       //
//   - UserId                               字符串 | 用户ID号码       //
// [私有属性] : -无-                                                  //
// [构造函数] : -无-                                                  //
// [公开函数] : -无-                                                  //
// [私有函数] : -无-                                                  //
////////////////////////////////////////////////////////////////////////

type LoginData struct {
	UserId string `name:"id" value:"" usage:"登录id"`
}

////////////////////////////////////////////////////////////////////////
//                        主要函数(runner)实现                        //
//--------------------------------------------------------------------//
// 实现函数:                                                          //
//        LoginRun(params commands.Params)       登录处理             //
//--------------------------------------------------------------------//
// 使用的内部的参数结构体(Params.Info属性对应的结构体) ： LoginData   //
////////////////////////////////////////////////////////////////////////

func LoginRun(params commands.Params) error {
	data := params.Info.(*LoginData)

	//判断参数中的userid
	if data.UserId == "" {
		//检查默认值
		if len(params.Args) > 0 {
			data.UserId = params.Args[0]
		}
	}

	//获取连接
	connect := params.Bundle[Connect].(net.Conn)
	//判断userId是否存在
	username, err := GlobalUserService.GetName(data.UserId)
	log.Print("检查用户ID(",data.UserId,")在线情况：",GlobalOnlineService.OnlineCheckByUserId(data.UserId))
	//如果存在返回登陆成功
	if err == nil  {
		if !GlobalOnlineService.OnlineCheckByUserId(data.UserId) {
			//加入在线列表
			GlobalOnlineService.Add(data.UserId, username, &connect)
			//写入output
			output := params.Bundle[Output].(io.Writer)
			fmt.Fprintln(output, data.UserId)
			log.Print("用户ID(",data.UserId,")登陆成功")
		} else {
			return errors.New("用户ID(" + data.UserId + ")已在线")
		}
	} else {
		log.Print("用户ID(",data.UserId,")登陆失败")
		return err
	}

	return nil
}
