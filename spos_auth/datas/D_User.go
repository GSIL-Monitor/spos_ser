// user
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"spos_auth/claims"
	"time"
)

//func AddUser() error {

//}
// AuthIdentity auth identity session model
type UserInfo struct {
	gorm.Model
	Basic
	SignLogs
}

// Basic basic information about auth identity
type Basic struct {
	UID string `gorm:"type:varchar(30);column:uid;unique_index"`

	UserName          string `gorm:"type:varchar(50);unique;not null"`
	EncryptedPassword string `gorm:"type:varchar(150);not null"`
	PhoneNum          string `gorm:"type:varchar(50)"`          // phone num
	MailAddr          string `gorm:"type:varchar(50);not null"` //email address
	ConfirmedAt       *time.Time
}

func (p *UserInfo) AddUser() error {

	ret := gdb.NewRecord(*p) // => returns `true` as primary key is blank
	fmt.Printf("AddUser:ret=%v\n", ret)

	ret1 := gdb.Create(p)
	fmt.Printf("AddUser:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func (p *UserInfo) GetUserByName(name string) error {

	ret := gdb.Where("user_name = ?", name).Find(p)

	return ret.Error
}

func (p *UserInfo) GetUserById(id string) error {

	ret := gdb.Where("uid = ?", id).Find(p)

	return ret.Error
}

func (p *UserInfo) UpdateUserConfirmedAt() error {
	t := time.Now()
	p.ConfirmedAt = &t
	fmt.Printf("UpdateUserConfirmedAt: uid=%s confirmed_at update\n", p.UID)
	ret := gdb.Model(p).Where("uid = ?", p.UID).Update("confirmed_at", p.ConfirmedAt)
	fmt.Printf("UpdateUserConfirmedAt:ret.Error=%v\n", ret.Error)

	return ret.Error
}

func IsExistUser(name string) bool {
	var user UserInfo

	gdb.Where("user_name = ?", name).Find(&user)

	fmt.Printf("IsExistUser: user=%v\n", user)

	//	fmt.Printf("IsExistUser: UserName=%s\n", user.UserName)
	if len(user.UserName) > 0 {
		return true
	}

	return false
}

func IsActivateUser(uid string) bool {
	var user UserInfo
	gdb.Where("uid = ?", uid).Find(&user)
	fmt.Printf("IsActivateUser: user=%v\n", user)
	//	fmt.Printf("IsActivateUser:  confirmed_at=%v\n", user.ConfirmedAt)
	t := user.ConfirmedAt
	if t == nil || t.IsZero() {
		fmt.Printf("IsActivateUser: User not activate confirm\n")
		return false
	}

	fmt.Printf("IsActivateUser: User confirmed\n")
	return true
}

// ToClaims convert to auth Claims
func (basic Basic) ToClaims() *claims.Claims {
	claims := claims.Claims{}
	//	claims.Provider = basic.Provider
	//	claims.Id = basic.UID
	claims.UserID = basic.UID
	return &claims
}

func GetAllUser() ([]UserInfo, error) {
	var users []UserInfo
	ret := gdb.Find(&users)

	return users, ret.Error
}
