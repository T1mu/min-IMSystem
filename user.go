package main

import (
	"fmt"
	"net"
	"strings"
)

// User 用户结构体，包含用户信息、Net.Conn连接和私人管道
type User struct {
	Name   string      // 用户名称
	Addr   string      // 用户地址
	Conn   net.Conn    // 用户对应的连接
	C      chan string // 用户对应的消息通道
	Server *Server
}

// NewUser 创建User对象，初始化对象。
func NewUser(conn net.Conn, server *Server) *User {
	user := &User{
		Name:   conn.RemoteAddr().String(),
		Addr:   conn.RemoteAddr().String(),
		Conn:   conn,
		C:      make(chan string),
		Server: server,
	}
	go user.listenUserChan()
	return user
}

// Online 用户上线功能
func (u *User) Online() {
	u.Server.UserMapMutex.Lock()
	u.Server.UserMap[u.Addr] = u
	u.Server.UserMapMutex.Unlock()
	u.Server.BroadCast(u, "我已上线")
}

// Offline 用户下线
func (u *User) Offline() {
	u.Server.UserMapMutex.Lock()
	delete(u.Server.UserMap, u.Addr)
	u.Server.UserMapMutex.Unlock()
	u.Server.BroadCast(u, "我已下线")
}

// SendMsg 向指定用户的客户端写消息
func (u *User) SendMsg(msg string) {
	u.Conn.Write([]byte(msg + "\n"))
}

// DoMsg 用户处理消息
func (u *User) DoMsg(msg string) {
	// 统计在线用户
	if msg == "who" {
		u.Server.UserMapMutex.Lock()
		for _, user := range u.Server.UserMap {
			u.SendMsg(fmt.Sprintf("[%s]:%s 在线\n", user.Addr, user.Name))
		}
		u.SendMsg(fmt.Sprintf("共计%d人\n", len(u.Server.UserMap)))
		u.Server.UserMapMutex.Unlock()
		// 通过rename方法修改用户名
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		u.Name = strings.Split(msg, "|")[1]
		u.SendMsg(fmt.Sprintf("新用户名为:[%s]\n", u.Name))

	} else if len(msg) > 5 && msg[:3] == "to|" {
		// 判断对方姓名是否存在，若存在则发送
		name := strings.Split(msg, "|")[1]
		var desUser *User
		u.Server.UserMapMutex.Lock()
		for _, user := range u.Server.UserMap {
			if user.Name == name {
				desUser = user
			}
		}
		u.Server.UserMapMutex.Unlock()
		if desUser == nil {
			u.SendMsg(fmt.Sprintf("未查找到[%s]用户\n", name))
		} else {
			cont := strings.Split(msg, "|")[2]
			desUser.SendMsg(fmt.Sprintf("[%s]%s:%s\n", u.Addr, u.Name, cont))
			u.SendMsg("发送成功！")
		}
		// 通过广播管道到私人管道实现群聊功能
	} else {
		u.Server.BroadCast(u, msg)
	}
}

// 监听私人管道，若有消息则回馈给服务端
func (u *User) listenUserChan() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}
