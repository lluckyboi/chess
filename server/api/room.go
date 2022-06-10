package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)
//断开ws时的handler
func Rec(roomid string, conn *websocket.Conn) {
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("coon.ReadMessage() err:", err)
			//如果是客户端发送，再close一遍
			conn.Close()
			delete(room, roomid)
			delete(roomcount,roomid)
			return
		}
	}
}

func enterroom(c *gin.Context){
	conn, err := up.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.JSON(200, gin.H{
			"play": false,
			"enter": false,
			"info":  "failed",
		})
		return
	}
	roomid:=c.PostForm("roomId")
	//userId:=c.MustGet("Id")

	//可以加入可以观战
		mplock.Lock()
		room[roomid] = conn
		roomcount[roomid]++
		mplock.Unlock()
		go Rec(roomid, conn)

		if roomcount[roomid]!=2{
		//进去下棋
		c.JSON(200,gin.H{
			"play":true,
			"enter":true,
			"info":roomid,
			"num":roomcount[roomid],
		})
	}else{
		//观战去
		c.JSON(200,gin.H{
			"play":false,
			"enter":true,
			"info":"人数已满",
			"num":roomcount[roomid],
		})
	}
}

func checkroom(c *gin.Context){
	roomId:=c.PostForm("roomId")
	if roomcount[roomId]>=2{
		c.JSON(200,gin.H{
			"status":true,
		})
	}else{

	}
}
