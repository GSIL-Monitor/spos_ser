// RoleObject
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type RoleObject struct {
	gorm.Model
	RID string `gorm:"type:varchar(30);column:rid;not null"`
	OID string `gorm:"type:varchar(30);column:oid;not null"`
}

func AddRoleObject(rid, oid string) error {
	info := RoleObject{RID: rid, OID: oid}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddRoleObject:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddRoleObject:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistRoleObject(rid, oid string) bool {

	var t RoleObject

	gdb.Where("rid = ? And oid = ?", rid, oid).Find(&t)

	fmt.Printf("IsExistRoleObject: t=%v\n", t)

	if len(t.RID) > 0 {
		return true
	}

	return false
}
