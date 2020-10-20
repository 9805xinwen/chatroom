package server

import (
	"bytes"
	"chatroom/server/cmds"
	"chatroom/server/services/message"
	"chatroom/server/services/users"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type ServerOpts struct {
	Name string
	Host string
	Port int
}

type Server struct {
	*ServerOpts
	listenTo string
	listener net.Listener
}

//设置默认值
func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Host == "" {
		newOpts.Host = "::"
	} else {
		newOpts.Host = opts.Host
	}
	if opts.Port == 0 {
		newOpts.Port = 3000
	} else {
		newOpts.Port = opts.Port
	}
	if opts.Name == "" {
		newOpts.Name = "ABC chatroom"
	} else {
		newOpts.Name = opts.Name
	}

	return &newOpts
}

//初始化全局表量，包括用户表、在线用户表、连接
func initSomething()  {
	cmds.GlobalUserService = users.NewRedisUser()
	cmds.GlobalOnlineService = users.NewMapOnline()
	cmds.GlobalMassageService = message.SimpleMessageService{}
}

//新建服务器
func NewServer(opts *ServerOpts) *Server {
	initSomething()

	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = net.JoinHostPort(opts.Host, strconv.Itoa(opts.Port))
	return s
}

//开始监听
func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", server.listenTo)
	if err != nil {
		return err
	}

	log.Printf("%s listening on %d", server.Name, server.Port)

	return server.Serve(listener)
}

//获取每一个请求的连接
func (server *Server) Serve(l net.Listener) error {
	server.listener = l
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		}
		go handle(tcpConn)
	}
}

func handle(conn net.Conn) {
	lineBuf := make([]byte, 1024)  //用于从conn里读取数据
	userIdBuf := bytes.NewBuffer(nil)

	var login bool  //是否已登陆
	var loginTimes = 3  //登录次数
	var userId string  //已登录的用户id

	for {
		n, err := conn.Read(lineBuf)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			break
		}

		line := string(lineBuf[:n])

		if !login {
			log.Printf("未认证主机[%s]：%s", conn.RemoteAddr(), line)
		} else {
			log.Printf("用户[ID:%s]: %s", userId, line)
		}

		cmdName := strings.ToLower(strings.SplitN(line, " ", 2)[0])
		cmdName = strings.TrimSpace(cmdName)

		if !login {
			loginTimes --
			if cmdName != cmds.LoginCommandName {
				io.WriteString(conn, "请先登录\n")
				loginTimes ++
				continue
			} else {
				loginBudle := map[string]interface{}{
					cmds.Connect: conn,
					cmds.Output: userIdBuf,
				}

				if err := cmds.LoginCommand.Execute(line, loginBudle); err != nil {
					if !login && loginTimes <= 0 {
						io.WriteString(conn, "登录次数用尽，即将退出\n")
						//登录次数用尽
						conn.Close()
						break
					}
					//登陆失败
					msg := fmt.Sprintf("登录失败，您还有%d次机会\n", loginTimes)
					io.WriteString(conn, msg)
					continue
				} else  {
					//登陆成功
					io.WriteString(conn, "登录成功\n")
					login = true
					userId, _ = userIdBuf.ReadString('\n')
					userId = strings.Fields(userId)[0]
					continue
				}
			}
		}

		cmd, ok := cmds.CommandMap[cmdName]
		if !ok {
			io.WriteString(conn, "命令不存在\n")
			continue
		}

		bundle := map[string]interface{}{
			cmds.Connect: conn,
			cmds.UserId: userId,
		}
		if err := cmd.Execute(line, bundle); err != nil {
			io.WriteString(conn, err.Error())
		}
	}
}
