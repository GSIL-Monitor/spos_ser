// UserOrg
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserOrg struct {
	gorm.Model
	UID  string `gorm:"type:varchar(30);column:uid;not null"`
	ORID string `gorm:"type:varchar(30);column:orid;not null"`
	DEID string `gorm:"type:varchar(30);column:deid;not null"`
}

func AddUserOrg(uid, orid, deid string) error {
	info := UserOrg{UID: uid, ORID: orid, DEID: deid}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddUserOrgHandler:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddUserOrgHandler:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistUserOrg(uid, orid, deid string) bool {

	var t UserOrg

	gdb.Where(" uid = ? And orid = ? And deid = ?", uid, orid, deid).Find(&t)

	fmt.Printf("IsExistUserOrg: t=%v\n", t)

	if len(t.UID) > 0 {
		return true
	}

	return false
}
