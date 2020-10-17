package server

import (
	"chatroom/common/util"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

type Auth interface {
	SignIn(string) string
	CheckId(string) (string, bool)
}

//一个基于redis的认证器
type SimpleAuth struct {
	db *redis.Client
}

func NewSimpleAuth() *SimpleAuth {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	auth := &SimpleAuth{rdb}
	return auth
}

//用户注册
func (auth *SimpleAuth) SignIn(name string) string {
	rrand := util.NewRandomUnique(auth.db, 10000, 99999)
	id := strconv.FormatInt(rrand.Next(), 10)
	hkey := fmt.Sprintf("USER:%s", id)
	if err := auth.db.HSet(hkey, "name", name).Err(); err != nil {
		panic(err)
	}
	return id
}

//用户登录认证
func (auth *SimpleAuth) CheckId(id string) (string, error) {
	hkey := fmt.Sprintf("USER:%s", id)
	name, err := auth.db.HGet(hkey, "name").Result()

	return name, err
}
