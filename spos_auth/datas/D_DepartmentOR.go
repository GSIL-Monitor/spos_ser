// DepartmentRole
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

type DepartmentOR struct {
	gorm.Model
	DEID string `gorm:"type:varchar(30);column:deid;not null"`
	OID  string `gorm:"type:varchar(30);column:oid;not null"`
}

func AddDepartmentOR(deid, oid string) error {
	info := DepartmentOR{OID: oid, DEID: deid}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddDepartmentRole:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddDepartmentRole:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistDepartmentOR(deid, oid string) bool {

	var t DepartmentOR

	gdb.Where("oid = ? And deid = ?", oid, deid).Find(&t)

	fmt.Printf("IsExistDepartmentOR: t=%v\n", t)

	if len(t.OID) > 0 {
		return true
	}

	return false
}
