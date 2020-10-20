package main

import "chatroom/server"

func main() {
	s := server.NewServer(&server.ServerOpts{})
	s.ListenAndServe()
}
