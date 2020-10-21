package users

import (
	"chatroom/common/util"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
)

//Users 一个接口，定义了操作总用户表的增加和查询方法
type Users interface {
	Register(name string) (id string, err error)  //注册
	GetName(id string) (name string, err error)  //根据id获取name
	GetId(name string) (id string, err error) //根据name获取id
}

//RedisUsers 以Redis存储用户信息，实现了User接口
type RedisUsers struct {
	db *redis.Client
}

//创建一个新的用户列表，redis地址指定为localhost:6379
func NewRedisUser() *RedisUsers {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	if _, err := rdb.Ping().Result(); err != nil {
		panic(err)
	}
	auth := &RedisUsers{rdb}
	return auth
}

func (user *RedisUsers) Register(name string) (id string, err error) {
	nameKey := fmt.Sprintf("USER:NAME:%s", name)

	if res, _ := user.db.Exists(nameKey).Result(); res != 0 {
		return "", errors.New("用户名已存在")
	}

	rrand := util.NewRandomUnique(user.db, 10000, 99999)
	id = strconv.FormatInt(rrand.NextWithKey("user"), 10)

	idKey := fmt.Sprintf("USER:ID:%s", id)

	pipe := user.db.TxPipeline()
	defer pipe.Close()

	pipe.HSetNX(idKey, "name", name)
	pipe.HSetNX(nameKey, "id", id)

	pipe.Exec()

	return id, nil
}


func (user *RedisUsers) GetName(id string) (string, error) {
	hkey := fmt.Sprintf("USER:ID:%s", id)
	name, err := user.db.HGet(hkey, "name").Result()

	return name, err
}

func (user *RedisUsers) GetId(name string) (string, error) {
	hkey := fmt.Sprintf("USER:NAME:%s", name)
	id, err := user.db.HGet(hkey, "id").Result()

	return id, err
}
