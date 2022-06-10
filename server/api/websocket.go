package api

import (
	"MyChess/client/chess"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	DelayTime         = time.Second / 20 //延迟
	ReceiveDataLength = 258*4            //接受数据长度
)


var wg sync.WaitGroup //进程锁
var up = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} //websocket协议升级结构体

var MsgCh = make(chan chess.PositionStruct, 50)       //信息广播通道
var mplock sync.Mutex                      			  //map锁
var mlock sync.Mutex                        		  //map锁
var room = make(map[string]*websocket.Conn)           //房间map
var roomcount = make(map[string]int)

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
	var mes chess.PositionStruct
	data := make([]byte, ReceiveDataLength)
	//持续读取数据
	for {
		_, err := conn.Read(data)
		if err != nil {
			//出口
			log.Println(err)
			return
		}
		//防止错误数据影响后续读入
		var check []byte
		check = data[0:ReceiveDataLength]
		errr := json.Unmarshal(check, &mes)
		if errr != nil {
			log.Println(errr)
			continue
		}
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
				if Msg.RoomId == i {
					j.WriteJSON(Msg)
				}
			}
			mplock.Unlock()
		}
	}
}
