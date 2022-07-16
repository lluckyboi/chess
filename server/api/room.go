package api

import (
	"MyChess/server/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type EnterRoomResp struct {
	Play  bool   `json:"play"`
	Enter bool   `json:"enter"`
	Info  string `json:"info"`
	Num   int    `json:"num"`
}

//断开ws时的handler
func Rec(roomid string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			tool.Logger.Info("Success..",
				zap.String("info", "coon.ReadMessage() err"),
				zap.Error(err))
			//如果是客户端发送，再close一遍
			conn.Close()
			delete(room, conn)
			roomcount[roomid]--
			return
		}
	}
}

func enterroom(c *gin.Context) {
	conn, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		tool.Logger.Info("Success..",
			zap.String("info", "enter room err"),
			zap.Error(err))
		c.JSON(200, gin.H{
			"play":  false,
			"enter": false,
			"info":  "failed",
			"num":   0,
		})
		return
	}

	roomid := c.Param("roomId")
	//userId:=c.MustGet("Id")

	//可以加入可以观战
	mplock.Lock()
	room[conn] = roomid
	roomcount[roomid] += 1
	mplock.Unlock()

	resp := EnterRoomResp{}
	if roomcount[roomid] <= 2 {
		resp = EnterRoomResp{
			Play:  true,
			Enter: true,
			Info:  "success",
			Num:   roomcount[roomid],
		}
	} else {
		resp = EnterRoomResp{
			Play:  false,
			Enter: true,
			Info:  "success",
			Num:   roomcount[roomid],
		}
	}

	fmt.Println(resp)
	conn.WriteJSON(resp)
	go Rec(roomid, conn)
}

func checkroom(c *gin.Context) {
	roomId := c.PostForm("roomId")
	if roomcount[roomId] >= 2 {
		c.JSON(200, gin.H{
			"status": true,
		})
	} else {
		c.JSON(200, gin.H{
			"status": false,
		})
	}
	fmt.Println(roomcount[roomId])
}
