package chess

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func UpdateBoard(g *Game) {
	for {
		_, message, err := Conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		json.Unmarshal(message, &g.singlePosition)
	}
}

func UpdateOtherBoard(g *Game) {
	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Println("net.Dail err", err)
		return
	}
	defer conn.Close()
	// 主动写数据给服务器
	ps := g.singlePosition
	str, err := json.Marshal(ps)
	if err != nil {
		log.Println("updateboard err", err)
		return
	}
	conn.Write(str)
	return
}
