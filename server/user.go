package server

import (
	"chatroom/common/util"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

type User interface {
	Register(name string) (id string, err error)
	GetName(id string) (name string, err error)
	GetId(name string) (id string, err error)
}

//一个基于redis的认证器
type RedisUser struct {
	db *redis.Client
}

func NewRedisUser() *RedisUser {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := rdb.Ping().Result(); err != nil {
		panic(err)
	}
	auth := &RedisUser{rdb}
	return auth
}

//用户注册, 用户表为id，name的一一映射
func (user *RedisUser) Register(name string) (id string, err error) {
	rrand := util.NewRandomUnique(user.db, 10000, 99999)
	id = strconv.FormatInt(rrand.NextWithKey("user"), 10)

	idKey := fmt.Sprintf("USER:ID:%s", id)
	nameKey := fmt.Sprintf("USER:NAME:%s", name)

	pipe := user.db.TxPipeline()
	defer pipe.Close()

	if _, err := pipe.Do("HSETNX", idKey, "name", name).Result(); err != nil {
		return "", errors.New("id重复")
	}
	if _, err := pipe.Do("HSETNX", nameKey, "id", id).Result(); err != nil  {
		return "", errors.New("name重复")
	}

	if _, err := pipe.Exec(); err != nil {
		return "", err
	}

	return id, nil
}


func (user *RedisUser) GetName(id string) (string, error) {
	hkey := fmt.Sprintf("USER:ID:%s", id)
	name, err := user.db.HGet(hkey, "name").Result()

	return name, err
}

func (user *RedisUser) GetId(name string) (string, error) {
	hkey := fmt.Sprintf("USER:NAME:%s", name)
	id, err := user.db.HGet(hkey, "id").Result()

	return id, err
}