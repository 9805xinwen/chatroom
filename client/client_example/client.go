package main

import (
	"bufio"
	"chatroom/client/client_com"
	"fmt"
	"log"
	"net"
	"os"
)

var writeStr, readStr = make([]byte, 1024), make([]byte, 1024)
var (
	chanQuit = make(chan bool, 0)
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalln(err)
		log.Println("Check connection settings")
		os.Exit(0)
	}

	defer conn.Close()

	fmt.Printf("%#v$请登录(命令 id)\n", conn.RemoteAddr().String())

	go handleSend(conn)
	go handleReceive(conn)
	<-chanQuit

}

func handleSend(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err0 := reader.ReadString('\n')
		if err0 != nil {
			fmt.Println("读取输入错误！")
			break
		}
		//把字符串分为【命令 参数】两部分
		cmd, _ := client_com.SpliteTwo(line)
		//判断输入的命令是否在map中
		if _, ok := client_com.Commands[cmd]; !ok {
			fmt.Println("输入命令不在列表中！请重新输入:")
			continue
		}
		//执行命令
		client_com.Commands[cmd].Execute()

		//给服务器发送命令
		_, err := conn.Write([]byte(line))
		if err != nil {
			fmt.Println("发送信息失败，err:", err)
			return
		}
	}
}

func handleReceive(conn net.Conn) {
	for {
		//接收服务端的消息
		n, err := conn.Read(readStr)
		if err != nil {
			fmt.Println("接收消息失败，err:", err)
		}
		if n > 0 {
			//将接收到的消息输出
			msg := string(readStr[:n])
			fmt.Println(msg)
		}

	}
}
