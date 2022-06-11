package chess

import (
	"encoding/json"
	"log"
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

}
