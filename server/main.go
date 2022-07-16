package main

import (
	"MyChess/server/api"
	"MyChess/server/dao"
	"MyChess/server/tool"
)

func main() {
	//先连接数据库
	dao.RUNDB()
	//初始化logger
	tool.InitLogger()
	defer tool.Logger.Sync()
	//接受机器数据的服务
	sev := api.NewServer("0.0.0.0", 8084)
	//广播ws
	go api.WsBroadcast()
	//启动tcp服务
	go sev.Start()
	//启动引擎
	api.RUNENGINE()
}
