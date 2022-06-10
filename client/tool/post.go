package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)
type MailResp struct {
	Code    int  `json:"code"`
	Info   string`json:"info"`
}

type LoginResp struct {
	Code    int  `json:"code"`
	UserName string`json:"username"`
}

type GetTokenResp struct {
	Code int `json:"code"`
	Info string `json:"info"`
	Token string `json:"token"`
}

type EnterRoomResp struct {
	Play bool `json:"play"`
	Enter bool `json:"enter"`
	Info string `json:"info"`
	Num   int  `json:"num"`
}

type CheckRoomCountResp struct {
	Status bool `json:"status"`
}

func GetMailAc(mail string)MailResp{
	resp, err := http.PostForm("http://localhost:9921/getmailac",
		url.Values{"UserMail": {mail}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg:=MailResp{}
	json.Unmarshal(body,&msg)
	return msg
}

func Login(mail,name,accode string)LoginResp{
	resp, err := http.PostForm("http://localhost:9921/login",
		url.Values{"UserMail": {mail},"UserName":{name},"accessCode":{accode}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg:=LoginResp{}
	json.Unmarshal(body,&msg)
	return msg
}

func GetToken(mail,name,accode string)GetTokenResp{
	resp, err := http.PostForm("http://localhost:9921/auth",
		url.Values{"UserMail": {mail},"UserName":{name},"accessCode":{accode}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg:=GetTokenResp{}
	json.Unmarshal(body,&msg)
	return msg
}

func EnterRoom(roomId string)EnterRoomResp{
	resp, err := http.PostForm("http://localhost:9921/enterroom",
		url.Values{"roomId":{roomId}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg:=EnterRoomResp{}
	json.Unmarshal(body,&msg)
	return msg
}

func CheckRoomCount(roomId string)CheckRoomCountResp{
	resp, err := http.PostForm("http://localhost:9921/checkroomcount",
		url.Values{"roomId":{roomId}})
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	msg:=CheckRoomCountResp{}
	json.Unmarshal(body,&msg)
	return msg
}