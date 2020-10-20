package server

import (
	"chatroom/server/services/message"
	"chatroom/server/services/users"
	"log"
	"net"
	"strconv"
)

var (
	UserList   *users.RedisUsers
	OnlineList *users.MapOnline
	Conns      []*net.Conn
	SMS        message.SimpleMessageService
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
func initSomething() {
	UserList = users.NewRedisUser()
	OnlineList = users.NewMapOnline()
	Conns = []*net.Conn{}
	SMS = message.SimpleMessageService{}
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

		Conns = append(Conns, &tcpConn)
	}
}
