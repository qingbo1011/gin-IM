# Gin+Websocket+MangoDB实现的IM系统

[基于 WebSocket + MongoDB 的IM即时聊天Demo](https://blog.csdn.net/weixin_45304503/article/details/121787022)

一个小demo。主要是熟悉一下WebSocket和MangoDB。（目前只支持单聊和文字，只是提前熟悉一下，为后面的游戏项目打个基础）

- `MySQL` ：存储用户基本信息
- `MongoDB` ：存放用户聊天信息
- `Redis` ：存储处理过期信息

实现功能：

- [x] 登录注册+JWT
- [x] 唯一token有效（同一个用户在同一个设备只能有一个有效token，如果产生新的token了，旧的token就算没有过期也会无效）
- [x] 单聊
- [ ] 群聊
- [x] 文字
- [x] emoji
- [ ] 语音，图片（表情包）
- [ ] 敏感词检查

其他功能有兴趣可以再自行研究。语音可图片可以考虑使用OSS之类的。这里只是一个demo，就只先完成一部分功能了。

## 配置文件

首先在conf下创建`config.ini`文件如下：

```ini
[service]
HttpPort = :8080

[mysql]
MysqlHost = 127.0.0.1
MysqlPort = 3308
MysqlUser = root
MysqlPassWord = 123456
MysqlName = chat_demo

[redis]
RedisAddr = 127.0.0.1:6379
RedisPassWord =
RedisDbName = 1

[mongo]
MangoAuthMechanism = SCRAM-SHA-1
MangoUser = root
MangoPassword = 1234
MangoHosts = 127.0.0.1:27017
MangoConnectTimeout = 10
MangoMaxPoolSize = 20
MangoMinPoolSize = 20
```

## WebSocket

参考笔记：[WebSocket编程](https://www.qingbo1011.top/2022/04/25/Golang%E8%BF%9B%E9%98%B601%20%E7%BD%91%E7%BB%9C%E7%BC%96%E7%A8%8B/#websocket%E7%BC%96%E7%A8%8B)

**[WebSocket 是什么原理？为什么可以实现持久连接？](https://www.zhihu.com/question/20215561/answer/40316953)**

- **WebSocket是一种在单个TCP连接上进行全双工通信的协议**
- WebSocket使得客户端和服务器之间的数据交换变得更加简单，**允许服务端主动向客户端推送数据**
- 在WebSocket API中，**浏览器和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输**

## 唯一token有效

- 第一次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzNzAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTczMDF9.ipiIDgAdTwrv8EX45y0UD6wy0fOOdzhIDysyB8kJais
- 第二次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzODAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTc0MDF9.3ZDrBr0FaFKpcicJpNkvEVCd8UdEQp079mg4fr2jBcc

通过测试可以发现，两个token都是有效的。只有到了指定的日期后，token才会失效。这显然是不合理的。在同一个设备上重新登陆后，

> 



