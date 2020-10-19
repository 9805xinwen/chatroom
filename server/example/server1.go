package main

import "chatroom/server"

func main() {
	redisUser := server.NewRedisUser()
	opts := &server.ServerOpts{User: redisUser}
	s := server.NewServer(opts)
	s.ListenAndServe()
}
