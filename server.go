package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// Server 服务端结构体，包含普通的IP、端口地址，新增用户哈希表数据、数据操作锁和公共管道
type Server struct {
	IP           string
	Port         int
	UserMap      map[string]*User
	UserMapMutex sync.RWMutex
	GlobChan     chan string
}

// NewServer 创建Sever方法，相当于初始化列表，并将Sever指针化
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:       ip,
		Port:     port,
		UserMap:  make(map[string]*User),
		GlobChan: make(chan string),
	}
	return server
}

// BroadCast send global message to GlobChan
func (p *Server) BroadCast(user *User, msg string) {
	msg = "[" + user.Addr + "]: " + user.Name + ": " + msg
	p.GlobChan <- msg
}

// Handler 连接成功的回调函数，具体功能：若服务端与用户端建立连接。
// 首先, 将用户数据放入用户数据表中，并将用户上线信息放入公共管道中。
func (p *Server) Handler(conn net.Conn) {
	user := NewUser(conn, p)
	// user模块 上线功能
	user.Online()
	// 提示有用户进入
	fmt.Sprintf("有新用户进入，其地址为%s", user.Addr)
	// 用户操作激活标志位
	alive := make(chan bool)
	// 接受客户端发送的消息
	go func() {
		for {
			buff := make([]byte, 4096)
			n, err := conn.Read(buff)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Err", err)
				return
			}
			msg := string(buff[:n-1])
			fmt.Println("服务端读到的内容:", msg)
			user.DoMsg(msg)
			alive <- true
		}
	}()
	for {
		select {
		// 触发激活
		// 所有 channel 表达式都会被求值，即超时函数也会被重新记时
		case <-alive:
			// 触发超时时限
		case <-time.After(time.Second * 60):
			// 关闭通道
			user.SendMsg("超时未操作，请重新连接\n")
			// 销毁资源
			close(user.C)
			// 关闭连接
			conn.Close()
			return
		}
	}
}

// listenGlobChan 监听管道
func (p *Server) listenGlobChan() {
	for {
		msg := <-p.GlobChan
		p.UserMapMutex.Lock()
		for _, i := range p.UserMap {
			i.C <- msg
		}
		p.UserMapMutex.Unlock()
	}
}

// Start 开始服务，分为监听IP、端口，若有则返回一个listner。
// 通过listener建立连接，返回conn。
// 再通过conn，处理相应Handler
func (p *Server) Start() {
	// 监听服务
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.IP, p.Port))
	if err != nil {
		fmt.Println("监听错误", err)
	}
	// 退出前关闭监听器
	defer listener.Close()
	// 监听通道
	go p.listenGlobChan()
	// 通过轮询监听器 建立连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("建立连接错误\n", err)
			continue
		}
		// 通过连接 做自己想做的事情 Handler
		go p.Handler(conn)
	}
	// 监听通道GlobChan
}
