// HeartService
package main

import (
	"encoding/json"
	"log"
	"spos_auth/datas"
)

const (
	USER_HEART_TIMEOUT = 6 //2 //mintue
)

type HeartResInfoList struct {
	ress []HeartResInfo `json:"heartres"`
}

type HeartResInfo struct {
	//	UniUserInfo
	Id       string `json:"id"`
	IsLogin  bool   `json:"islogin"` //1---login; 0----logout;
	IsOnline bool   `json:"isonline"`
	IsRecvOk bool   `json:"isrecvok"`
}

func GetHeartResInfo() string {
	//var users datas.UserInfo
	users, _ := datas.GetAllUser()
	//	log.Printf("GetHeartResInfo: users num=%d\n", len(users))

	var reslist HeartResInfoList
	var islogin, isonline, isrecvok bool
	for _, user := range users {
		//		log.Printf("GetHeartResInfo list: user id=%s; name=%s\n", user.UID, user.UserName)

		islogin, isonline, isrecvok = GetUserStatusFromWorkList(user.UID)
		reslist.ress = append(reslist.ress, HeartResInfo{Id: user.UID, IsLogin: islogin, IsOnline: isonline, IsRecvOk: isrecvok})

	}

	b, err := json.Marshal(reslist.ress)
	if err != nil {
		log.Println("json err:", err)
	}

	log.Printf("GetHeartResInfo list json=%s\n", string(b))

	return string(b)
}
