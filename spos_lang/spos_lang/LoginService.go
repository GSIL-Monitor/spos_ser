// LoginService
package main

import (
	"encoding/json"
	"log"
	"spos_auth/datas"
)

type LoginResInfoList struct {
	ress []LoginResInfo `json:"loginres"`
}

type LoginResInfo struct {
	//	UniUserInfo
	Id   string `json:"id"`
	Name string `json:"name"`
	Gid  string `json:"gid"`
	//	Mac      string `json:"mac"`
	IsLogin  bool `json:"islogin"` //1---login; 0----logout;
	IsOnline bool `json:"isonline"`
}

func GetLoginResInfo() string {
	//var users datas.UserInfo
	users, _ := datas.GetAllUser()
	//	log.Printf("GetLoginResInfo: users num=%d\n", len(users))

	var reslist LoginResInfoList
	var islogin, isonline bool
	for _, user := range users {
		//		log.Printf("GetLoginResInfo list: user id=%s; name=%s\n", user.UID, user.UserName)

		islogin, isonline, _ = GetUserStatusFromWorkList(user.UID)
		reslist.ress = append(reslist.ress, LoginResInfo{Id: user.UID, Name: user.UserName, IsLogin: islogin, IsOnline: isonline})

	}

	b, err := json.Marshal(reslist.ress)
	if err != nil {
		log.Println("json err:", err)
	}

	log.Printf("GetLoginResInfo list json=%s\n", string(b))

	return string(b)
}
