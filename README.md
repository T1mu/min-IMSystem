# min-IMSystem
min-IMSystem是一个由Go语言编写的**入门级**的项目，它整体只由三个go文件组成，分别是client.go, server.go 和 user.go。
## server.go
server.go表示服务端，由若干种方法构成。
### start方法
其中Start方法用于创建一个net.Conn的连接，创建一个Go程不断监听server对象的管道内容（listenGlobChan），创建一个Go程读取并处理客户端发送的消息（Handler）。
### listenGlobChan
监听GlobChan公共管道内容，若有内容，则将内容分发到每个用户的UserChan私人管道中。
### Handler
- 处理用户上下线事件（方法在user包中，在server中调用）
- 处理读取连接内容信息
- 监听用户是否超时
## user.go
user.go为服务端包含控制用户上下线、传输信息给用户等功能的模块
### online offline
在对server对象的用户HashMap哈希表加锁、解锁前提下，实现用户上线下线功能。
- 上线online：增加用户到哈希表中
- 下线offline：删除用户到哈希表中
### sendMsg
调用user的连接conn属性的Write方法，从而实现通知用户当前消息是否处理成功的功能。
### listenUserChan
监听用户私人管道，若管道中有新内容，则调用Write方法或sendMsg方法。
### DoMsg
用于处理服务端Read到的信息，具体为：判断用户写入数据（服务端读取到的数据）是否符合一定规则，若符合一定规则则按需求执行。
## client.go
客户端，取代unix下nc命令
### Run
Run方法打印提示用户可执行的功能，若用户选择某种功能，比如：公聊模式，则客户端模拟向连接中写入公聊模式格式，之后，服务端通过读连接中的数据，将数据处理后传入User对象的DoMsg方法中。