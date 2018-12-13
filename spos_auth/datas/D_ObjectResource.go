// ObjectResource
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type ObjectResource struct {
	gorm.Model
	OID      string `gorm:"type:varchar(30);column:oid;unique_index;not null"`
	Name     string `gorm:"type:varchar(50);unique;not null"`
	Content  string `gorm:"type:varchar(100);not null"`
	Action   string `gorm:"type:varchar(30);not null"`
	Describe string `gorm:"type:varchar(100)"`
}

func AddObjectResource(name, content, action, describe string) error {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	info := ObjectResource{OID: "O" + strconv.FormatInt(timestamp, 10), Name: name, Content: content, Action: action, Describe: describe}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddObjectResourceHandler:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddObjectResourceHandler:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func GetObjectResource(oid string) *ObjectResource {

	var t ObjectResource

	gdb.Where("oid = ?", oid).Find(&t)

	fmt.Printf("IsExistObjectResource: t=%v\n", t)
	//	fmt.Printf("IsExistObjectResource: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return &t
	}

	return nil
}

func IsExistObjectResource(name string) bool {

	var t ObjectResource

	gdb.Where("name = ?", name).Find(&t)

	fmt.Printf("IsExistObjectResource: t=%v\n", t)
	//	fmt.Printf("IsExistObjectResource: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return true
	}

	return false
}
