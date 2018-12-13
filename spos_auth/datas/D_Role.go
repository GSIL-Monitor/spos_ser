// Role
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Role struct {
	gorm.Model
	RID      string `gorm:"type:varchar(30);column:rid;unique_index;not null"`
	Name     string `gorm:"type:varchar(50);not null"`
	Describe string `gorm:"type:varchar(100)"`
}

func AddRole(name, describe string) error {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	info := Role{RID: "R" + strconv.FormatInt(timestamp, 10), Name: name, Describe: describe}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddRole:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddRole:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistRole(name string) bool {
	var t Role

	gdb.Where("name = ?", name).Find(&t)

	fmt.Printf("IsExistRole: t=%v\n", t)

	fmt.Printf("IsExistRole: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return true
	}

	return false
}
