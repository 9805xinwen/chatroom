package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"strconv"
)

var ErrServerClosed = errors.New("chatroom: Server closed")

type ServerOpts struct {
	User     User
	Name     string
	Hostname string
	Port     int

	WelcomeMessage string
}

type Server struct {
	*ServerOpts
	listenTo  string
	listener  net.Listener
	tlsConfig *tls.Config
	ctx       context.Context
	cancel    context.CancelFunc
	feats     string
}

func serverOptsWithDefaults(opts *ServerOpts) *ServerOpts {
	var newOpts ServerOpts
	if opts == nil {
		opts = &ServerOpts{}
	}
	if opts.Hostname == "" {
		newOpts.Hostname = "::"
	} else {
		newOpts.Hostname = opts.Hostname
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

	if opts.WelcomeMessage == "" {
		newOpts.WelcomeMessage = defaultWelcomeMessage
	} else {
		newOpts.WelcomeMessage = opts.WelcomeMessage
	}

	if opts.User != nil {
		newOpts.User = opts.User
	}

	return &newOpts
}

func NewServer(opts *ServerOpts) *Server {
	opts = serverOptsWithDefaults(opts)
	s := new(Server)
	s.ServerOpts = opts
	s.listenTo = net.JoinHostPort(opts.Hostname, strconv.Itoa(opts.Port))
	return s
}

func (server *Server) newConn(tcpConn net.Conn) *Conn {
	c := new(Conn)
	c.conn = tcpConn
	c.controlReader = bufio.NewReader(tcpConn)
	c.controlWriter = bufio.NewWriter(tcpConn)
	c.user = server.User
	c.server = server
	c.sessionID = newSessionID()

	return c
}

func (server *Server) ListenAndServe() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", server.listenTo)
	if err != nil {
		return err
	}

	sessionID := ""
	log.Printf("%s %s listening on %d", sessionID, server.Name, server.Port)

	return server.Serve(listener)
}

// Serve accepts connections on a given net.Listener and handles each
// request in a new goroutine.
//
func (server *Server) Serve(l net.Listener) error {
	server.listener = l
	server.ctx, server.cancel = context.WithCancel(context.Background())
	sessionID := ""
	for {
		tcpConn, err := server.listener.Accept()
		if err != nil {
			select {
			case <-server.ctx.Done():
				return ErrServerClosed
			default:
			}
			log.Printf("%s listening error: %v", sessionID, err)
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		}

		conn := server.newConn(tcpConn)
		go conn.Serve()
	}
}

// Shutdown will gracefully stop a server. Already connected clients will retain their connections
func (server *Server) Shutdown() error {
	if server.cancel != nil {
		server.cancel()
	}
	if server.listener != nil {
		return server.listener.Close()
	}
	// server wasnt even started
	return nil
}
