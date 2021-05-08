package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flg        int
}

// NewClient 创建客户端
func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flg:        999, //flg若为默认值0，则Run方法第一次执行即退出循环
	}
	// 连接服务器，获取返回的conn
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("拨号过程出现错误\n", err)
		return nil
	}
	// 设置客户端的conn
	client.conn = conn
	return client
}

// menu 显示打印供用户操作文本，并返回成功或失败
func (c *Client) menu() bool {
	// 用户输入数字变量
	var flg int
	// 打印用户操作提示
	fmt.Println("1 - 公聊模式")
	fmt.Println("2 - 私聊模式")
	fmt.Println("3 - 更新用户名")
	fmt.Println("0 - 退出")
	// 读取用户输入数字
	fmt.Scanln(&flg)

	// 判断数字合法性
	if flg >= 0 && flg < 5 {
		// 将操作模式赋予客户端对象flg属性
		c.flg = flg
		return true
	} else {
		fmt.Println(">>>>>>>>请输入合法数字")
		return false
	}
}

// publicChat 通过conn实现公共聊天
func (c *Client) publicChat() bool {
	var msg string
	fmt.Println(">>>>>>>>请输入公聊内容")
	fmt.Scanln(&msg)
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("数据流写入出错", err)
		return true
	}
	return false
}

// privateChat 通过conn写入实现私人聊天
func (c *Client) privateChat() bool {
	var name, msg string
	fmt.Println("请输入私人聊天对方名称")
	fmt.Scanln(&name)
	fmt.Println("请输入私人聊天内容")
	fmt.Scanln(&msg)
	msg = "to|" + name + "|" + msg + "\n"
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("数据流写入出错", err)
		return false
	}
	return true
}

// updateName 通过conn写入rename|新名称数据，模拟用户修改名称
func (c *Client) updateName() bool {
	fmt.Println(">>>>>>>>请输入用户名")
	fmt.Scanln(&c.Name)
	sendMsg := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("数据流写入出错:", err)
		return false
	}
	return true
}

// dealResponse 阻塞状态读取conn消息
func (c *Client) dealResponse() {
	io.Copy(os.Stdout, c.conn) //注意：此方法永久阻塞，不会只执行一次
	// 其含义是将c.conn的数据读到标准输出上

}

// Run 根据不同的模式处理不同的业务
func (c *Client) Run() {
	for c.flg != 0 {
		// 若用户输入一直为错，则一直调用menu
		// 注意：每次调用客户端对象的menu方法都会打印操作提示文本，并且设置c.flg的值
		for c.menu() != true {
		}
		switch c.flg {
		case 1:	// 公共聊天
			c.publicChat()
		case 2:	// 私聊
			c.privateChat()
		case 3:
			c.updateName()
		}
	}
}

// 使用flag库，对命令行参数进行解析
var serverIp string
var serverPort int

// init 通过flag库绑定命令行参数、默认值和提示
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器地址")
	flag.IntVar(&serverPort, "p", 7327, "设置服务器端口号")
}

func main() {
	// 命令行解析
	flag.Parse()
	// 通过NewClient创建客户端
	client := NewClient(serverIp, serverPort)
	// 判断用户是否创建成功，通过nil变量
	if client == nil {
		fmt.Println(">>>>>>>>连接服务器失败")
		return
	}
	// 提示用户创建成功
	fmt.Println(">>>>>>>>连接服务器成功")
	// 启动一个go程读
	go client.dealResponse()
	// 客户端处理业务
	client.Run()
}
