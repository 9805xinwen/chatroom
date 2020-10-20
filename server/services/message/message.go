package message

import (
	"fmt"
	"net"
)

//Message 封装了发件方、收件方、消息内容
type Massage struct {
	FromUser string
	ToUser string
	FromConn *net.Conn
	ToConn *net.Conn
	Content string
}

//MassageService 提供消息的传递服务
type MassageService interface {
	Send(massage Massage) error
}

type SimpleMessageService struct {}

func (sms SimpleMessageService) Send(message Massage) error {
	//log.Printf("%s向%s发送消息：%s", message.FromUser, message.ToUser, message.Content)
	fmt.Fprintf(*message.ToConn, "%s: %s", message.FromUser, message.Content)
	return nil
}