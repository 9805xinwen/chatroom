package server

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	defaultWelcomeMessage = "Welcome to the chatroom"
)

//一个用户一条连接
type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	user          User
	server        *Server
	sessionID     string
	reqUserId     string
	userId        string
	userName      string
	closed        bool
}

//生成唯一id标识连接
func newSessionID() string {
	hash := sha256.New()
	_, err := io.CopyN(hash, rand.Reader, 50)
	if err != nil {
		return "????????????????????"
	}
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr[0:20]
}

//为单条连接serve
func (conn *Conn) Serve() {
	log.Print(conn.sessionID, "Connection Established")
	conn.writeMessage(defaultWelcomeMessage)

	for {
		line, err := conn.controlReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Print(conn.sessionID, fmt.Sprint("read error:", err))
			}

			break
		}
		conn.receiveLine(line)
		// QUIT command closes connection, break to avoid error on reading from
		// closed socket
		if conn.closed == true {
			break
		}
	}
	conn.Close()
	log.Print(conn.sessionID, "Connection Terminated")
}

//处理请求
func (conn *Conn) receiveLine(line string) {
	command, param := conn.parseLine(line)
	log.Printf("%s %s %s", conn.sessionID, command, param)
	cmdObj := commands[strings.ToUpper(command)]
	if cmdObj == nil {
		return
	}
	//if cmdObj.RequireParam() && param == "" {
	//
	//} else
	if cmdObj.RequireAuth() && conn.userId == "" {

	} else {
		cmdObj.Execute(conn, param)
	}
}

//解析请求
func (conn *Conn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

//向连接中写入信息
func (conn *Conn) writeMessage(message string) (wrote int, err error) {
	log.Printf("发送给：%s %s", conn.sessionID, message)
	line := fmt.Sprintf("%s\r\n", message)
	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}

//关闭连接
func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
}
