// UniSposUsers
package main

import (
	"container/list"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	USER_LOGOUT = 0
	USER_LOGIN  = 1
)
const (
	USER_OFFLINE = 0
	USER_ONLINE  = 1
)

type UniUserInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Gid  string `json:"gid"`
	Mac  string `json:"mac"`
}

type UniWorkUser struct {
	UniUserInfo
	IsLogin       bool   `json:"islogin"` //1---login; 0----logout;
	LoginTime     string `json:"logintime"`
	LogoutTime    string `json:"logouttime"`
	IsOnline      bool
	IsRecvOk      bool
	LastHeartTime time.Time
	//	MsgCh         *chan UniMsg
}

type Queue struct {
	l *list.List
	m sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{l: list.New()}
}

func (q *Queue) PushBack(v interface{}) {
	if v == nil {
		return
	}
	q.m.Lock()
	defer q.m.Unlock()
	q.l.PushBack(v)
}

func (q *Queue) Front() *list.Element {
	q.m.Lock()
	defer q.m.Unlock()
	return q.l.Front()
}

func (q *Queue) Remove(e *list.Element) {
	if e == nil {
		return
	}
	q.m.Lock()
	defer q.m.Unlock()
	q.l.Remove(e)
}

func (q *Queue) Len() int {
	q.m.Lock()
	defer q.m.Unlock()
	return q.l.Len()
}

var (
	work_user_list = list.New()

	q_user = Queue{l: work_user_list}
)

func UpdateUserToWorkList(user UniWorkUser) bool {

	if len(user.Id) <= 0 {
		return false
	}

	log.Printf("UpdateUserToWorkList::user id=%s;mac=%s updateing\n", user.Id, user.Mac)

	if GetWorkUserByIdAndMac(user.Id, user.Mac) != nil {
		RemoveUserFromWorkList(user)
	}

	AddUserToWorkList(user)

	log.Printf("UpdateUserToWorkList::user work list ok. num=%d\n", q_user.Len())

	return true
}

func AddUserToWorkList(user UniWorkUser) bool {
	var c UniWorkUser
	var n *list.Element

	if len(user.Id) <= 0 {
		return false
	}

	log.Printf("AddUserToWorkList::user id=%s;mac=%s,work_user_list num=%d\n", user.Id, user.Mac, q_user.Len())

	for e := q_user.Front(); e != nil; e = n {

		n = e.Next()

		c = e.Value.(UniWorkUser)

		//		log.Printf("AddUserToWorkList::user id=%s;mac=%s;islogin=%v; isonline=%v\n", c.Id, c.Mac, c.IsLogin, c.IsOnline)

		if strings.EqualFold(c.Id, user.Id) && strings.EqualFold(c.Mac, user.Mac) {
			log.Printf("AddUserToWorkList::user is exist in work list\n")
			return false
		}

	}

	q_user.PushBack(user)

	log.Printf("AddUserToWorkList::user added list ok. num=%d\n", q_user.Len())

	return true
}

func GetWorkUserByIdAndMac(id, mac string) *UniWorkUser {

	var c UniWorkUser
	var n *list.Element

	log.Printf("GetWorkUserByIdAndMac::id=%s; mac=%s\n", id, mac)

	for e := q_user.Front(); e != nil; e = n {
		n = e.Next()

		c = e.Value.(UniWorkUser)
		//		log.Printf("GetWorkUserByIdAndMac::c.Id=%s; c.Mac=%s \n", c.Id, c.Mac)
		if strings.EqualFold(id, c.Id) && strings.EqualFold(mac, c.Mac) {
			log.Printf("GetWorkUserByIdAndMac::get id=%s mac=%s user success\n", id, mac)
			return &c
		}

	}

	log.Printf("GetWorkUserByIdAndMac::get id=%s user failed\n", id)

	return nil
}

func GetWorkUserById(id string) []UniWorkUser {

	var us []UniWorkUser
	var c UniWorkUser
	var n *list.Element

	log.Printf("GetWorkUserById::id=%s\n", id)

	for e := q_user.Front(); e != nil; e = n {
		n = e.Next()

		c = e.Value.(UniWorkUser)
		//		log.Printf("GetWorkUserById::c.Id=%s\n", c.Id)
		if strings.EqualFold(id, c.Id) {
			log.Printf("GetWorkUserById::get id=%s user success\n", id)
			us = append(us, c)

		}

	}

	log.Printf("GetWorkUserById::get id=%s user num=%d\n", id, len(us))

	return us
}

func RemoveUserFromWorkList(u UniWorkUser) bool {

	var c UniWorkUser
	var n *list.Element
	log.Printf("RemoveUserFromWorkList user id=%s mac=%s deleteing..\n", u.Id, u.Mac)

	for e := q_user.Front(); e != nil; e = n {

		n = e.Next()

		c = e.Value.(UniWorkUser)

		if strings.EqualFold(c.Id, u.Id) && strings.EqualFold(c.Mac, u.Mac) {
			//			close(*c.MsgCh)
			//			c.MsgCh = nil
			q_user.Remove(e)
			log.Printf("RemoveUserFromWorkList user id=%s mac=%s removed ok. user list num=%d\n", c.Id, c.Mac, q_user.Len())
			return true
		}

	}
	log.Printf("RemoveUserFromWorkList user not exist and removed fail")
	return false
}

func GetUserStatusFromWorkList(id string) (bool, bool, bool) {
	var c UniWorkUser
	var n *list.Element

	islogin := false
	isonline := false
	isrecvok := false

	//	log.Printf("GetUserStatusFromWorkList  id=%s\n", id)

	for e := q_user.Front(); e != nil; e = n {

		n = e.Next()

		c = e.Value.(UniWorkUser)

		//	log.Printf("id=%s; c.Id=%s \n", id, c.Id)
		if strings.EqualFold(id, c.Id) {
			//			log.Printf("get id=%s mac=%s islogin=%v isonline=%v user Status success\n", id, c.Mac, c.IsLogin, c.IsOnline)
			if c.IsLogin {
				islogin = true
			}
			if c.IsOnline {
				isonline = true
			}
			if c.IsRecvOk {
				isrecvok = true
			}
		}

	}

	//	log.Printf("GetUserStatusFromWorkList id=%s user islogin=%v isonline=%v\n", id, islogin, isonline)

	return islogin, isonline, isrecvok
}

func GetUserMacFromWorkList(id string) string {
	var c UniWorkUser
	var n *list.Element

	log.Printf("GetUserStatusFromList  id=%s\n", id)

	for e := q_user.Front(); e != nil; e = n {

		n = e.Next()

		c = e.Value.(UniWorkUser)

		//		ulog.Ul.Debug("GET list c nw1=", c.Uconn.RemoteAddr().String())
		//		ulog.Ul.Debug("GET conn nw1=", conn.RemoteAddr().String())
		//	log.Printf("id=%s; c.Id=%s \n", id, c.Id)
		if strings.EqualFold(id, c.Id) {
			log.Printf("get id=%s user Mac success\n", id)
			return c.Mac
		}

	}

	log.Printf("get id=%s user Mac failed\n", id)

	return "null"
}

func WorkUserHeartTimeOutHandle() {
	var c UniWorkUser
	var n *list.Element

	log.Printf("Lang WorkUserHeartTimeOutHandle begin")

	for e := q_user.Front(); e != nil; e = n {
		n = e.Next()

		c = e.Value.(UniWorkUser)

		//		m, _ := time.ParseDuration("-1m")
		sub := time.Now().Sub(c.LastHeartTime)
		log.Printf("HeartTimeOutHandle::Id=%s mac=%s lastt=%v sub=%v IsOnline=%v\n",
			c.Id, c.Mac, c.LastHeartTime, sub, c.IsOnline)
		if int(sub.Seconds()) > USER_HEART_TIMEOUT*60 && c.IsOnline {
			user := GetWorkUserByIdAndMac(c.Id, c.Mac)
			if user != nil {
				user.IsOnline = false
				UpdateUserToWorkList(*user)
			}
		}

	}

	log.Printf("Lang WorkUserHeartTimeOutHandle end")
}
