/**
 *@Author:sario
 *Date:2022/7/16
 *@Desc:
 */
package tool

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

type MailResp struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

type LoginResp struct {
	Code     int    `json:"code"`
	UserName string `json:"username"`
}

type GetTokenResp struct {
	Code  int    `json:"code"`
	Info  string `json:"info"`
	Token string `json:"token"`
}

type EnterRoomResp struct {
	Play  bool   `json:"play"`
	Enter bool   `json:"enter"`
	Info  string `json:"info"`
	Num   int    `json:"num"`
}

type CheckRoomCountResp struct {
	Status bool `json:"status"`
}

type WinResp struct {
	WinCount int `json:"win_count"`
}

const addr = "http://39.106.81.229"
const port = "9924"

func GetMailAc(mail string) MailResp {
	resp, err := http.PostForm(addr+":"+port+"/getmailac",
		url.Values{"UserMail": {mail}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg := MailResp{}
	json.Unmarshal(body, &msg)
	return msg
}

func Login(mail, name, accode string) LoginResp {
	resp, err := http.PostForm(addr+":"+port+"/login",
		url.Values{"UserMail": {mail}, "UserName": {name}, "accessCode": {accode}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg := LoginResp{}
	json.Unmarshal(body, &msg)
	return msg
}

func GetToken(mail, name, accode string) GetTokenResp {
	resp, err := http.PostForm(addr+":"+port+"/auth",
		url.Values{"UserMail": {mail}, "UserName": {name}, "accessCode": {accode}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg := GetTokenResp{}
	json.Unmarshal(body, &msg)
	return msg
}

func EnterRoom(roomId string, token string) (EnterRoomResp, *websocket.Conn) {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer "+token)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "39.106.81.229:9924", Path: "/enterroom/" + roomId}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		log.Fatal("dial:", err)
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
	}
	resp := EnterRoomResp{}
	json.Unmarshal(message, &resp)

	log.Println(resp)
	return resp, c
}

func CheckRoomCount(roomId string) CheckRoomCountResp {
	resp, err := http.PostForm(addr+":"+port+"/checkroomcount",
		url.Values{"roomId": {roomId}})
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	msg := CheckRoomCountResp{}
	json.Unmarshal(body, &msg)
	return msg
}

func AddWinCount(token string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", addr+":"+port+"/addwin", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	defer resp.Body.Close()
}

func GetWinCount(token string) int {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", addr+":"+port+"/getwincount", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	resp, _ := client.Do(req)
	ws := WinResp{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	json.Unmarshal(body, &ws)
	defer resp.Body.Close()
	return ws.WinCount
}
