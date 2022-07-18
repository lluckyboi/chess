# 功能描述
![image](https://img.shields.io/badge/testes-66.7%25-green)

可以通过客户端，登陆注册 , 加入房间对战、观战

## 技术细节

棋盘通过**ebtien** 库绘制，键鼠交互

客户端调用服务器接口

邮箱注册登录

进入房间，人满2人自动开始，他人可进入观战

客户端发起**tcp请求**，发送棋盘数据给服务器，服务器**推送到websocke**t , 客户端再通过协程**监听websocket连接读取**棋盘数据，然后进行**实时渲染**

客户端负责逻辑处理，服务端仅负责收发数据，因此并发较高

**棋局回放**：服务器拿到棋盘数据根据用户信息入库，后续客户端访问接口，服务器定时从websocket推数据，服务器拿数据本地渲染回放（未做完）

## 使用方法

**client**文件夹单独放在一个项目中，可以多开几个项目测试，运行即可，按照提示输入邮箱和用户名，然后查看验证码，输入验证码，输入房间号（最好6位），若未满两人会进行等待，满两人自动开始，先进入的为红方

或者通过EXE文件双击运行（因为未知原因 会闪退到后台）

## 接口简述

```go
//发送邮箱验证码
r.POST("/getmailac", getmailac)
//登录注册
r.POST("/login", login)
//拿token
r.POST("/auth", authHandler)
//进入房间 进入升级为ws连接 然后进行初始化 本地维护room和roomcount
r.GET("/enterroom/:roomId", JWTAuthMiddleware(), enterroom)
//检查人数
r.POST("/checkroomcount", checkroom)
```

