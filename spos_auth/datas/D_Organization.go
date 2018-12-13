// Organization
package datas

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"strconv"
	"time"
)

type Organization struct {
	gorm.Model
	ORID       string `gorm:"type:varchar(30);column:orid;unique_index;not null"`
	ICID       string `gorm:"type:varchar(30);column:icid;not null"`
	Name       string `gorm:"type:varchar(50);not null"`
	AgentName  string `gorm:"type:varchar(50);not null"`
	AgentPhone string `gorm:"type:varchar(50);not null"`
	AgentMail  string `gorm:"type:varchar(50);not null"`
	Describe   string `gorm:"type:varchar(100)"`
}

func AddOrganization(icid, name, agentname, agentphone, agentmail, describe string) error {
	timestamp := time.Now().Unix()
	fmt.Println(timestamp)

	info := Organization{ORID: "ORG" + strconv.FormatInt(timestamp, 10), ICID: icid, Name: name,
		AgentName: agentname, AgentPhone: agentphone, AgentMail: agentmail, Describe: describe}

	ret := gdb.NewRecord(info) // => returns `true` as primary key is blank
	fmt.Printf("AddOrganization:ret=%v\n", ret)

	ret1 := gdb.Create(&info)
	fmt.Printf("AddOrganization:ret1.Error=%v\n", ret1.Error)

	return ret1.Error
}

func IsExistOrganization(name string) bool {
	var t Organization

	gdb.Where("name = ?", name).Find(&t)

	fmt.Printf("IsExistOrganization: t=%v\n", t)

	//	fmt.Printf("IsExistOrganization: Name=%s\n", t.Name)
	if len(t.Name) > 0 {
		return true
	}

	return false
}
