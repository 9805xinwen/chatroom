package users

import "net"

//Online 一个接口 定义了操作用户在线列表的方法
type Online interface {
	Add(userId, userName string, conn *net.Conn) bool  //添加某用户到在线列表，传入id、name、conn，返回是否添加成功
	OnlineCheckByUserId(userId string) bool  //检查某用户是否在线，传入id，返回是否在线
	OnlineCheckByUserName(userName string) bool  //检查某用户是否在线，传入name，返回是否在线
	QueryConnByUserId(userId string) *net.Conn  //查询指定id的连接，传入id，返回Conn
	QueryConnByUserName(userName string) *net.Conn  //查询指定name的连接，传入name，返回Conn
	Delete(userId string) bool  //删除在线列表中的某用户，传入id，返回是否删除成功
}

//MapOnline 使用map实现了Online接口
type MapOnline struct {
	mapId map[string]*net.Conn
	mapName map[string]*net.Conn
}

//创建一个新的用户在线列表
func NewMapOnline() *MapOnline {
	mapId := make(map[string]*net.Conn)
	mapName := make(map[string]*net.Conn)

	return &MapOnline{
		mapId: mapId,
		mapName: mapName,
	}
}

func (online *MapOnline) Add(usersId ,userName string, conn *net.Conn) bool {
	if _, ok := online.mapId[usersId]; !ok {
		return false
	}
	online.mapId[usersId] = conn
	online.mapId[userName] = conn
	return true
}

func (online *MapOnline) OnlineCheckByUserId(userId string) bool {
	if _, ok := online.mapId[userId]; !ok {
		return false
	}
	return true
}

func (online *MapOnline) OnlineCheckByUserName(userName string) bool {
	if _, ok := online.mapId[userName]; !ok {
		return false
	}
	return true
}

func (online *MapOnline) QueryConnByUserId(userId string) *net.Conn {
	conn, _ := online.mapId[userId]
	return conn
}

func (online *MapOnline) QueryConnByUserName(userName string) *net.Conn {
	conn, _ := online.mapId[userName]
	return conn
}

func (online *MapOnline) Delete(userId string) bool {
	delete(online.mapId, userId)
	return true
}

