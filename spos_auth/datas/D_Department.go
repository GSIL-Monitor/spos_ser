// Department
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Department struct {
	gorm.Model
	DEID     string `gorm:"type:varchar(30);column:deid;unique_index;not null"`
	ORID     string `gorm:"type:varchar(30);column:orid;not null"`
	Name     string `gorm:"type:varchar(50);not null"`
	Describe string `gorm:"type:varchar(100)"`
}

func AddDepartment(orid, name, describe string) error {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	info := Department{DEID: "D" + strconv.FormatInt(timestamp, 10), ORID: orid, Name: name, Describe: describe}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddDepartmentHandler:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddDepartmentHandler:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistDepartment(name string) bool {
	var t Department

	gdb.Where("name = ?", name).Find(&t)

	fmt.Printf("IsExistDepartment: t=%v\n", t)

	//	fmt.Printf("IsExistDepartment: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return true
	}

	return false
}
