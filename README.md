# Gin+Websocket+MangoDB实现的IM系统

[基于 WebSocket + MongoDB 的IM即时聊天Demo](https://blog.csdn.net/weixin_45304503/article/details/121787022)

一个小demo。主要是熟悉一下WebSocket和MangoDB。（目前只支持单聊和文字，只是提前熟悉一下，为后面的游戏项目打个基础）

- `MySQL` ：存储用户基本信息
- `MongoDB` ：存放用户聊天信息
- `Redis` ：存储处理过期信息

关于这个项目的用户登录注册以及JWT相关内容，跟gin-memos项目中一样。

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









