package main

import (
	"MyChess/server/api"
	"MyChess/server/dao"
)

func main()  {
	//先连接数据库
	dao.RUNDB()
	//接受机器数据的服务
	sev := api.NewServer("0.0.0.0", 8082)
	//广播ws
	go api.WsBroadcast()
	//启动tcp服务
	go sev.Start()
	//启动引擎
	api.RUNENGINE()
}
