/*
随机数生成器
用途：每次运行产生指定范围内的唯一存在的数字
用法：第一步，调用NewRandomUnique()工厂方法产生生成器r，此方法需要一个redis实例作为参数传入，同时传入min，max，
	第二步，调用r.Next()或r.Next(key string)返回一个int64
*/

package util

import (
	"github.com/go-redis/redis"
	"math/rand"
	"time"
)

type RandomUnique struct {
	Seed int64

	Min int64
	Max int64

	DB *redis.Client
}

func NewRandomUnique(db *redis.Client, min, max int64) *RandomUnique {
	r := RandomUnique{
		Seed: time.Now().UnixNano(),
		Min:  min,
		Max:  max,
		DB:   db,
	}

	return &r
}

func (r RandomUnique) Next() (res int64) {
	if r.DB.SCard("random:unique").Val() == (r.Max - r.Min) {
		panic("当前范围内的数字都已被使用")
	}
	rand.Seed(r.Seed)
	for {
		res = rand.Int63n(r.Max-r.Min) + r.Min
		if !r.DB.SIsMember("random:unique", res).Val() {
			r.DB.SAdd("random:unique", res)
			break
		}
	}
	return
}

func (r RandomUnique) NextWithKey(key string) (res int64) {
	if r.DB.SCard(key).Val() == (r.Max - r.Min) {
		panic("当前范围内的数字都已被使用")
	}
	rand.Seed(r.Seed)
	for {
		res = rand.Int63n(r.Max-r.Min) + r.Min
		if !r.DB.SIsMember(key, res).Val() {
			r.DB.SAdd(key, res)
			break
		}
	}
	return
}


/*

func main() {
	defer trace("rand")()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       1,  // use default DB
	})

	r := NewRandomUnique(rdb, 0, 1000)

	for i:=0; i<1000; i++ {
		r.Next()
	}
}

func trace(msg string) func() {
	start := time.Now()
	log.Printf("enter %s", msg)
	return func() {
		log.Printf("exit %s (%s)", msg, time.Since(start))
	}
}

*/
