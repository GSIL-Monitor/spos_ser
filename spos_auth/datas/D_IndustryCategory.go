// IndustryCategory
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type IndustryCategory struct {
	gorm.Model
	ICID     string `gorm:"type:varchar(30);column:icid;unique_index;not null"`
	Name     string `gorm:"type:varchar(50);not null"` // IndustryCategory Name
	Describe string `gorm:"type:varchar(100)"`
}

func AddIndustryCategory(name, describe string) error {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	info := IndustryCategory{ICID: "I" + strconv.FormatInt(timestamp, 10), Name: name, Describe: describe}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddIndustryCategory:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddIndustryCategory:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistIndustryCategory(name string) bool {
	var t IndustryCategory

	gdb.Where("name = ?", name).Find(&t)

	fmt.Printf("IsExistIndustryCategory: t=%v\n", t)

	//	fmt.Printf("IsExistIndustryCategory: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return true
	}

	return false
}
