package main

import (
	"chatroom/server"
	"chatroom/server/services"
)

func main() {
	redisUser := services.NewRedisUser()
	opts := &server.ServerOpts{User: redisUser}
	s := server.NewServer(opts)
	s.ListenAndServe()
}
