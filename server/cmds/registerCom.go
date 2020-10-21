package cmds

import (
	"chatroom/common/commands"
	"fmt"
	"io"
	"reflect"
)

////////////////////////////////////////////////////////////////////////
//                            Register 命令定义                       //
//--------------------------------------------------------------------//
// [命令名称] : Register                                              //
// [命令参数] :                                                       //
//              -username                 [默认] 想要注册的用户名     //
//--------------------------------------------------------------------//
// 使用案例:                                                          //
// register -username USERNAME                                        //
// register Username                                                  //
////////////////////////////////////////////////////////////////////////

const RegisterCommandName string = "register"

var RegisterCommand commands.Command = commands.CreateDefaultCommand(RegisterCommandName, reflect.TypeOf(RegisterData{}), RegisterRun)

////////////////////////////////////////////////////////////////////////
//                        主要命令参数结构体定义                      //
//--------------------------------------------------------------------//
// 定义结构体:                                                        //
//        RegisterData                          退出数据结构体        //
//--------------------------------------------------------------------//
// RegisterData                                                       //
// [公开属性] :                                                       //
//              - UserName                 字符串 | 想要注册的用户名  //
// [私有属性] : -无-                                                  //
// [构造函数] : -无-                                                  //
// [公开函数] : -无-                                                  //
// [私有函数] : -无-                                                  //
////////////////////////////////////////////////////////////////////////

type RegisterData struct{
	UserName string `name:"username" value:"" usage:"用户名"`
}

////////////////////////////////////////////////////////////////////////
//                        主要函数(runner)实现                        //
//--------------------------------------------------------------------//
// 实现函数:                                                          //
//        RegisterRun(params commands.Params)        退出处理         //
//--------------------------------------------------------------------//
// 使用的内部的参数结构体(Params.Info属性对应的结构体) ： RegisterData//
////////////////////////////////////////////////////////////////////////

func RegisterRun(params commands.Params) error {
	//获取数据
	data := params.Info.(*RegisterData)
	//如果用户没用按 register -username USERNAME 的方式来写
	//则默认按照 register USERNAME 的方式读取
	if data.UserName == "" {
		if len(params.Args) > 0 {
			data.UserName = params.Args[0]
		}
	}

	//读取id和错误
	id,err := GlobalUserService.Register(data.UserName)
	if err != nil {
		return err
	} else {
		output := params.Bundle[Output].(io.Writer)
		fmt.Fprintln(output,id)
	}
	return nil
}
