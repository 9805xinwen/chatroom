package message

import (
	"fmt"
	"net"
)

//Message 封装了发件方、收件方、消息内容
type Massage struct {
	FromUser string
	ToUser string
	FromConn    *net.Conn
	ToConn      *net.Conn
	Content string
}

//MassageService 提供消息的传递服务
type MassageService interface {
	Send(massage Massage) error
}

type SimpleMessageService struct {}

func (sms SimpleMessageService) Send(message Massage) error {
	to := *message.ToConn
	//do something

	msg := "hi" //最终发给客户端的消息
	fmt.Fprint(to, msg)
	//io.WriteString(to, msg)
	//io.Write([]byte(msg))
	return nil
}