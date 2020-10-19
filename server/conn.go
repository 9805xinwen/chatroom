package server

import (
	"bufio"
	"chatroom/server/cmds"
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

type Conn struct {
	conn          net.Conn
	controlReader *bufio.Reader
	controlWriter *bufio.Writer
	auth          Auth
	server        *Server
	sessionID     string
	reqUser       string
	user          string
	closed        bool
}

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

func (conn *Conn) Serve() {
	log.Print(conn.sessionID, "Connection Established")

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

func (conn *Conn) receiveLine(line string) {
	command, param := conn.parseLine(line)
	log.Printf("%s %s %s", conn.sessionID, command, param)
	cmdObj := cmds.commands[strings.ToUpper(command)]
	if cmdObj == nil {
		return
	}
	cmdObj.Execute(conn, param)
}

func (conn *Conn) parseLine(line string) (string, string) {
	params := strings.SplitN(strings.Trim(line, "\r\n"), " ", 2)
	if len(params) == 1 {
		return params[0], ""
	}
	return params[0], strings.TrimSpace(params[1])
}

// writeMessage will send a standard FTP response back to the client.
func (conn *Conn) writeMessage(message string) (wrote int, err error) {
	log.Print("%s %s %s", conn.sessionID, message)
	line := fmt.Sprintf("%s\r\n", message)
	wrote, err = conn.controlWriter.WriteString(line)
	conn.controlWriter.Flush()
	return
}

func (conn *Conn) Close() {
	conn.conn.Close()
	conn.closed = true
}
