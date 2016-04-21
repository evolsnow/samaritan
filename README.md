# **project samaritan**

[![Build Status](https://api.travis-ci.org/evolsnow/samaritan.svg?branch=master)](https://travis-ci.org/evolsnow/samaritan)
[![GoDoc](https://godoc.org/github.com/evolsnow/samaritan?status.png)](https://godoc.org/github.com/evolsnow/samaritan)

**samaritan** 为APP(Godo)的后台部分, 以Golang编写。

### 项目介绍
与另外两位同学的2016年毕业论文题为《基于iOS端的个人日程管理与团队管理工具的研究与实现》，

samaritan为其实现部分。此处仅简单介绍后台部分的一些架构及技术细节。

[项目文档](https://samaritan.gitbooks.io/app-doc/content/)

接口与文档细节正在逐步更新。

初版，以后可能会出于安全考虑而转为私有。

### 架构介绍
- WEB框架：挑选一些Idiomatic组件，手动搭建web框架。为求轻便与更深入地了解golang，故项目未采用已有的（过度设计）框架。

- 组织结构：并没有严格按照MVC结构，且VIEW层此时不需要。

- 数据存储：redis + mysql混合。redis用于热数据存取，mysql此处仅用来备份数据，方便以后统计分析等操作。

- 缓存设计：redis缓存连接池 + LRU缓存。前者用来缓存一些生存周期短的内容，如验证码之类; 高速LRU缓存主要用来存取用户的jwt计算结果，降低cpu负载，加快用户授权信息验证的过程。

- 安全考虑： 1.严格的jwt鉴权; 2.id混淆，避免URI猜测; 3.scrypt强加密密码

- 静态资源： 头像，图片等交由七牛云。

- RPC: 采用基于ProtoBuf的gRPC。消息推送(websocket/apns), 邮件发送等功能直接通过rpc实现, 形成微服务。非必须, 兴趣而已。

- TLS：虽然golang本身支持TLS，但实测下来效率不如nginx加解密HTTPS流量+golang处理HTTP请求的模式(与服务器应该有关系，待测)，所以直接交由nginx处理TLS，并开启HTTP/2支持。

### 引用库说明
golang的一大特点便是fork优秀的开源库再针对自己的项目进行优化使用，所以samaritan也是这么搭建而成的，真心感谢开源作者们。

- [negroni](https://github.com/codegangsta/negroni): Martini作者停止维护Martini后着手的新的web中间件管理项目，fork之后加了HTTPS支持。

- [binding](https://github.com/mholt/binding): HTTP请求内容与定义的结构体绑定中间件，规范请求参数，方便调试。fork后修改了错误内容的输出模式。

- [jwt-go](https://github.com/dgrijalva/jwt-go): [JSON Web Tokens](http://self-issued.info/docs/draft-jones-json-web-token.html)中间件的golang实现，用户请求信息鉴权。

- [httprouter](https://github.com/julienschmidt/httprouter): 简洁高效的路由管理，fork后增加了上下文参数存取的功能，方便在中间件及handler间传递。

- [logrus](https://github.com/Sirupsen/logrus): 更方便的日志输出，按照等级输出在debug时很有必要。

- [scrypt](https://golang.org/x/crypto/scrypt): 更为安全的加密方式，存储密码时用。

- [redigo](https://github.com/garyburd/redigo): redis连接池操作，项目中主要配合Lua脚本使用，减少连接次数。



### License
[The MIT License (MIT)](https://raw.githubusercontent.com/evolsnow/samaritan/master/LICENSE)
