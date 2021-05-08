# min-IMSystem
min-IMSystem是一个由Go语言编写的**入门级**的项目，它整体只由三个go文件组成，分别是client.go, server.go 和 user.go。
## server.go
server.go由若干种方法构成。
### start方法
其中Start方法用于创建一个net.Conn的连接，创建一个Go程不断监听server对象的管道内容（listenGlobChan），创建一个Go程读取并处理客户端发送的消息（Handler）。
### listenGlobChan
监听公共管道内容，若有内容，则将内容分发到每个客户端的私人管道中。
### Handler
- 处理用户上下线事件（方法在user包中，在server中调用）
- 处理读取连接内容信息
- 监听用户是否超时
## user.go
todo
## client.go
todo