# Gin+Websocket+MangoDB实现的IM系统

[基于 WebSocket + MongoDB 的IM即时聊天Demo](https://blog.csdn.net/weixin_45304503/article/details/121787022)

一个小demo。主要是熟悉一下WebSocket和MangoDB。（目前只支持单聊和文字，只是提前熟悉一下，为后面的游戏项目打个基础）

- `MySQL` ：存储用户基本信息
- `MongoDB` ：存储用户聊天信息
- `Redis` ：存储处理过期信息；存储token信息，保证同一个User-Agent只有唯一一个token有效。

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
MangoDBName = gin-IM
MangoConnectTimeout = 10
MangoMaxPoolSize = 20
MangoMinPoolSize = 20
```

## 唯一token有效

- 第一次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzNzAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTczMDF9.ipiIDgAdTwrv8EX45y0UD6wy0fOOdzhIDysyB8kJais
- 第二次登录的token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU3OTgzODAxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc4OTc0MDF9.3ZDrBr0FaFKpcicJpNkvEVCd8UdEQp079mg4fr2jBcc

通过测试可以发现，两个token都是有效的。只有到了指定的日期后，token才会失效。这显然是不合理的。所以我在这里的处理是：**引入Redis**。

具体逻辑：使用redis中的hash结构。`key`为user的唯一标识uid；`filed`为该user的User-Agent，表示是哪一个设备（同一个设备只能有1个token有效）；`value`存储该user的唯一有效token。

结构如下：

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220716184841.png)

- 在`service/user.go`中的`UserLogin()`增加逻辑：登录成功后签发token时，确定`uid`和`User-Agent`。直接`Rdb.HSet()`即可
  - 因为`Rdb.HSet()`的处理逻辑是，如果存在就更新value；不存在就新建。
- 在`middleware/jwt.go`中增加判定逻辑。通过`uid`和`User-Agent`（从解析token中包含的相关信息claims中获取uid），查出redis中的token。判定携带的token是否和redis中的token一样。如果不一样说明是旧的token，直接`c.Abort()`然后`return`。

> 这只是我自己的一个想法，如果以后发现更好的解决方案，会继续更新的。

更新后：

- 第一次token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU4MDU0NzUxLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc5NjgzNTF9.-FvhHHpJokeigiSJOUkTWaQ4ytsYDZcxaTklPLzJGR4
- 第二次token：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsInVzZXJuYW1lIjoidG9tIiwiZXhwIjoxNjU4MDU0NzgwLCJpc3MiOiJnaW4tSU0iLCJuYmYiOjE2NTc5NjgzODB9.uBkmCpbTfEbr3fBiMQ26XrxOQc-hl6H5jvS_3BfW-2o

可以发现使用第一次token去请求会403：

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220716184738.png)

> 据说微信就是这样做的（跟群友讨论的）：
>
> - 这不就是提掉线吗？
> - 登录后 将以前此用户的token删除掉即可
> - 如果想多设备登录 就加入设备就可以 当前token和用户id，设备绑定
> - 微信就是这样做的
>
> 跟大家讨论，感觉基本都是基于redis缓存token的，踢掉用户也是这么干的。
>
> 不过感觉这样就跟jwt的无状态背道而驰了，回到了session。如果以后有更优雅更好的方式，会再记录的。

## WebSocket

参考笔记：[WebSocket编程](https://www.qingbo1011.top/2022/04/25/Golang%E8%BF%9B%E9%98%B601%20%E7%BD%91%E7%BB%9C%E7%BC%96%E7%A8%8B/#websocket%E7%BC%96%E7%A8%8B)

**[WebSocket 是什么原理？为什么可以实现持久连接？](https://www.zhihu.com/question/20215561/answer/40316953)**

- **WebSocket是一种在单个TCP连接上进行全双工通信的协议**
- WebSocket使得客户端和服务器之间的数据交换变得更加简单，**允许服务端主动向客户端推送数据**
- 在WebSocket API中，**浏览器和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输**

关于websocket，这里为了方便我急没有使用jwt鉴权。在api里写了一个test接口，去测试token的唯一有效。然后关于websocket的代码中，`replayMsg`的code我这里为了方便就用的`http`状态码。这显然是不合理的。应该需要提前做好约定，然后自己封装一些关于websocket的状态码。

postman中测试websocket接口：**[Postman Now Supports WebSocket APIs](https://blog.postman.com/postman-supports-websocket-apis/)**

> postman9版本在win10上安装出错解决方案：Postman安装失败： [Installation has failed Failed to extract installer](https://blog.csdn.net/zhouyingge1104/article/details/119359357)
>

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220718155110.gif)

### 简单对聊

这里就看一下postman响应数据和Redis、MangoDB中存储到的数据即可：

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220718225228.png)![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220718225239.png)

![](https://img-qingbo.oss-cn-beijing.aliyuncs.com/img/20220718225258.png)

### Postman WebSocket JWT鉴权的坑

本来想在Postman中测试Websocket的鉴权的，但是发现websocket连接中获取到的User-Agent为空，需要自己在header中指定：

















