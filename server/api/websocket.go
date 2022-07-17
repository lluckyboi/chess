package api

import (
	"MyChess/server/tool"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	DelayTime         = time.Second / 20 //延迟
	ReceiveDataLength = 589              //接受数据长度
)

type PositionStruct struct {
	SdPlayer    int      `json:"SdPlayer"`    // 轮到谁走，0=红方，1=黑方
	UcpcSquares [256]int `json:"UcpcSquares"` // 棋盘上的棋子
	RoomId      string   `json:"RoomId"`      //棋手ID
}

var wg sync.WaitGroup //进程锁
var up = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} //websocket协议升级结构体

//type RoomAndLock struct {
//	MsgCh			chan[500] PositionStruct			//信息广播通道
//	RoomCountLock	sync.Mutex
//	RoomLock		sync.Mutex
//	Room 			map[*websocket.Conn]string			//房间map
//	RoomCount		map[string]int						//房间人数
//}

var MsgCh = make(chan PositionStruct, 500)  //信息广播通道
var mplock sync.Mutex                       //map锁
var mlock sync.Mutex                        //map锁
var room = make(map[*websocket.Conn]string) //房间map
var roomcount = make(map[string]int)        //房间人数

type Server struct {
	Ip   string
	Port int
} //服务器

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

// Start 启动
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net listen err:", err)
	}
	defer listener.Close()
	//持续监听端口
	for {
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("net Accept err:", err)
		}
		go this.broaddata(con)
	}
}

func (this *Server) broaddata(conn net.Conn) {
	log.Println("tcp 连接成功")
	var mes PositionStruct
	data := make([]byte, ReceiveDataLength)
	//持续读取数据
	for {
		n, err := conn.Read(data)
		if err != nil {
			//出口
			log.Println(err)
			return
		}
		var check []byte
		check = data[0:n]
		log.Println(n)
		errr := json.Unmarshal(check, &mes)
		if errr != nil {
			log.Println(errr)
			continue
		}
		log.Println(mes)
		MsgCh <- mes
	}
}

// WsBroadcast 广播websocket
func WsBroadcast() {
	for {
		select {
		case Msg := <-MsgCh:
			mplock.Lock()
			for i, j := range room {
				if Msg.RoomId == j {
					err := i.WriteJSON(Msg)
					if err != nil {
						tool.Logger.Info("ws board err",
							zap.Error(err))
					}
				}
			}
			mplock.Unlock()
		}
	}
}

//func test(c *gin.Context){
//	ps:=PositionStruct{
//		SdPlayer:    0,
//		UcpcSquares: [256]int{0,1,2,2,3,1,23,5,45,},
//		RoomId:      "12365",
//	}
//	c.JSON(200,gin.H{
//		"ps":ps,
//	})
//}
