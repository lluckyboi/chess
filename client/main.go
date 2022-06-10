package main

import (
	"MyChess/client/chess"
	"MyChess/client/tool"
	"fmt"
	"time"
)



func main() {
	fmt.Println("please login: input you mail and name")
	var mail,name,accode string
	//输入用户名邮箱
	fmt.Scanf("%s%s",&mail,&name)

	//获取邮箱验证码
	msg:=tool.GetMailAc(mail)
	fmt.Println(msg.Info)
	if msg.Code!=200{
		return
	}

	//输入验证码
	fmt.Println("input you mail access code:")
	fmt.Scanf("%s",&accode)

	//用户登录
	lmsg:=tool.Login(mail,name,accode)
	if lmsg.Code!=200{
		return
	}
	fmt.Println(lmsg.UserName+"登录注册成功")

	//拿token
	gmsg:=tool.GetToken(mail,name,accode)
	if gmsg.Code!=2000{
		return
	}

	//输入房间号
	var roomId string
	fmt.Println("please input roomId:")
	fmt.Scanf("%s",&roomId)

	//加入房间
	emsg:=tool.EnterRoom(roomId)
	if emsg.Enter{
		fmt.Println("成功进入房间！")
		if emsg.Play{
			fmt.Println("成功加入对局！")
		}else{
			fmt.Println("当前仅能观战！")
		}
	}else{
		fmt.Println("进入房间失败！")
		return
	}

	//预处理资源
	tool.FileToByte("./resource","./chess")

	//如果成功加入对战
	if emsg.Enter&&emsg.Play{
		//轮询，人满即开
		for{
			cmsg:=tool.CheckRoomCount(roomId)
			if cmsg.Status==true{
				break
			}
			time.Sleep(time.Second)
		}

		//启动游戏 先进入的为红方
		chess.NewGame(emsg.Num-1)
	}
}