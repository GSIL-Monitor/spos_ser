// UserRole
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type UserRole struct {
	gorm.Model
	UID string `gorm:"type:varchar(30);column:uid;not null"`
	RID string `gorm:"type:varchar(30);column:rid;not null"`
}

func AddUserRole(uid, rid string) error {
	info := UserRole{RID: rid, UID: uid}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddUserRoleHandler:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddUserRoleHandler:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistUserRole(rid, uid string) bool {

	var t UserRole

	gdb.Where("rid = ? And uid = ?", rid, uid).Find(&t)

	fmt.Printf("IsExistUserRole: t=%v\n", t)

	if len(t.RID) > 0 {
		return true
	}

	return false
}
