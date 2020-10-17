package main

import (
	"chatroom/client"
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)


var writeStr, readStr = make([]byte, 1024), make([]byte, 1024)
var (
	chanQuit = make(chan bool,0)
)

//把字符串分为【命令 参数】两部分
func spliteTwo(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}else{
		return params[0], strings.TrimSpace(params[1])
	}
}



func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Fatalln(err)
		log.Println("Check connection settings")
		os.Exit(0)
	}

	defer conn.Close()

	fmt.Printf("%#v$请登录(命令 id)\n", conn.RemoteAddr().String())

	//reader := bufio.NewReader(os.Stdin)
	//
	////判断输入字符串是否符合格式
	//for {
	//	line, _ := reader.ReadString('\n')
	//	//把字符串分为【命令 参数】两部分
	//	cmd, param := spliteTwo(line)
	//	//判断输入的命令是否在map中
	//	if _, ok := cm.Commands[cmd]; !ok {
	//		break
	//	} else {
	//		//判断该命令是否需要参数，且后面的参数格式是否正确
	//		//正确与否都继续输入
	//		cm.Commands[cmd].CommandFormat(param, cm.Commands[cmd].RequireParam())
	//		//执行命令
	//		//cm.Commands[cmd].Execute()
	//
	//		//给服务器发送命令
	//		_, err = conn.Write([]byte(line))
	//		if err != nil {
	//			fmt.Println("发送信息失败，err:", err)
	//			return
	//		}
	//		//收消息
	//		var buf [1024]byte
	//		n,err := conn.Read(buf[:])
	//		if err != nil{
	//			fmt.Println("接收消息失败，err:",err)
	//			return
	//		}
	//		fmt.Println("收到消息：",string(buf[:n]))
	//	}
	//}
	go handleSend(conn)
	go handleReceive(conn)
	<-chanQuit

}

func handleSend(conn net.Conn){
	reader := bufio.NewReader(os.Stdin)
	for {
		line, _ := reader.ReadString('\n')
		//把字符串分为【命令 参数】两部分
		cmd, param := spliteTwo(line)
		//判断输入的命令是否在map中
		if _, ok := client.Commands[cmd]; !ok {
			fmt.Println("输入命令不在列表中！请重新输入:")
			continue
		} else {
			//判断该命令是否需要参数，且后面的参数格式是否正确
			//正确与否都继续输入
			client.Commands[cmd].CommandFormat(param, client.Commands[cmd].RequireParam())
			//执行命令
			//cm.Commands[cmd].Execute()

			//给服务器发送命令
			_, err := conn.Write([]byte(line))
			if err != nil {
				fmt.Println("发送信息失败，err:", err)
				return
			}
		}
	}
}

func handleReceive(conn net.Conn){
	for{
		//接收服务端的消息
		n, err := conn.Read(readStr)
		if err != nil {
			fmt.Println("接收消息失败，err:",err)
		}
		if n > 0{
			//将接收到的消息输出
			msg := string(readStr[:n])
			fmt.Println(msg)
		}

	}
}