package datas

import (
	"errors"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"

	"github.com/jinzhu/gorm"

	"spos_auth/enforcer"
)

var (
	gdb *gorm.DB
)

func OpenUnDbAndInitCasbin() error {

	fmt.Println("OpenUnDb: WORDPRESS_DB_HOST:" + os.Getenv("WORDPRESS_DB_HOST"))
	fmt.Println("OpenUnDb: WORDPRESS_DB_PASSWORD:" + os.Getenv("WORDPRESS_DB_PASSWORD"))

	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(" + os.Getenv("WORDPRESS_DB_HOST") + ":3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(172.17.0.14:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	connect_str := "root:Unitone@2018" + "@tcp(192.168.20.107:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"

	fmt.Printf("OpenUnDb: connect_str=%s\n", connect_str)
	db, err := gorm.Open("mysql", connect_str)
	fmt.Printf("OpenUnDb: Open err=%v\n", err)

	if err != nil {
		return errors.New("db open error")
	} else {
		gdb = db
		gdb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&UserInfo{}, &Role{},
			&UserRole{}, &ObjectResource{}, &RoleObject{},
			&IndustryCategory{}, &Organization{}, &Department{},
			&DepartmentOR{}, &UserOrg{})

		fmt.Printf("OpenUnDb: UserRole create Foreign Key!\n")
		gdb.Model(&UserRole{}).AddForeignKey("uid", "user_infos(uid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserRole{}).AddForeignKey("rid", "roles(rid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: RoleObject create Foreign Key!\n")
		gdb.Model(&RoleObject{}).AddForeignKey("oid", "object_resources(oid)", "RESTRICT", "RESTRICT")
		gdb.Model(&RoleObject{}).AddForeignKey("rid", "roles(rid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: Organization create Foreign Key!\n")
		gdb.Model(&Organization{}).AddForeignKey("icid", "industry_categories(icid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: Department create Foreign Key!\n")
		gdb.Model(&Department{}).AddForeignKey("orid", "organizations(orid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: DepartmentRole create Foreign Key!\n")
		gdb.Model(&DepartmentOR{}).AddForeignKey("deid", "departments(deid)", "RESTRICT", "RESTRICT")
		gdb.Model(&DepartmentOR{}).AddForeignKey("oid", "object_resources(oid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: UserOrg create Foreign Key!\n")
		gdb.Model(&UserOrg{}).AddForeignKey("uid", "user_infos(uid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserOrg{}).AddForeignKey("orid", "organizations(orid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserOrg{}).AddForeignKey("deid", "departments(deid)", "RESTRICT", "RESTRICT")

		gdb.Model(&RoleObject{}).AddUniqueIndex("idx_roleobject_rid_oid", "rid", "oid")
		gdb.Model(&UserRole{}).AddUniqueIndex("idx_userrole_uid_rid", "uid", "rid")
		gdb.Model(&DepartmentOR{}).AddUniqueIndex("idx_departmentor_deid_oid", "deid", "oid")
		gdb.Model(&UserOrg{}).AddUniqueIndex("idx_userorg_uid_orid_deid", "uid", "orid", "deid")

	}

	fmt.Printf("OpenUnDb: db open ok!\n")

	enforcer.InitCasbin(connect_str)
	return nil
}

func OnlyOpenUnDb() error {

	fmt.Println("OnlyOpenUnDbOpenUnDb: WORDPRESS_DB_HOST:" + os.Getenv("WORDPRESS_DB_HOST"))
	fmt.Println("OnlyOpenUnDbOpenUnDb: WORDPRESS_DB_PASSWORD:" + os.Getenv("WORDPRESS_DB_PASSWORD"))

	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(" + os.Getenv("WORDPRESS_DB_HOST") + ":3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(" + os.Getenv("WORDPRESS_DB_HOST") + ":33061)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	//	connect_str := "root:" + os.Getenv("WORDPRESS_DB_PASSWORD") + "@tcp(172.17.0.14:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"
	//connect_str := "root:Unitone@2018" + "@tcp(192.168.20.107:3306)/uni_spos?charset=utf8&parseTime=True&loc=Local"

	fmt.Printf("OpenUnDb: connect_str=%s\n", connect_str)
	db, err := gorm.Open("mysql", connect_str)
	fmt.Printf("OpenUnDb: Open err=%v\n", err)

	if err != nil {
		return errors.New("db open error")
	} else {
		gdb = db
		gdb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(&UserInfo{}, &Role{},
			&UserRole{}, &ObjectResource{}, &RoleObject{},
			&IndustryCategory{}, &Organization{}, &Department{},
			&DepartmentOR{}, &UserOrg{})

		fmt.Printf("OpenUnDb: UserRole create Foreign Key!\n")
		gdb.Model(&UserRole{}).AddForeignKey("uid", "user_infos(uid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserRole{}).AddForeignKey("rid", "roles(rid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: RoleObject create Foreign Key!\n")
		gdb.Model(&RoleObject{}).AddForeignKey("oid", "object_resources(oid)", "RESTRICT", "RESTRICT")
		gdb.Model(&RoleObject{}).AddForeignKey("rid", "roles(rid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: Organization create Foreign Key!\n")
		gdb.Model(&Organization{}).AddForeignKey("icid", "industry_categories(icid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: Department create Foreign Key!\n")
		gdb.Model(&Department{}).AddForeignKey("orid", "organizations(orid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: DepartmentRole create Foreign Key!\n")
		gdb.Model(&DepartmentOR{}).AddForeignKey("deid", "departments(deid)", "RESTRICT", "RESTRICT")
		gdb.Model(&DepartmentOR{}).AddForeignKey("oid", "object_resources(oid)", "RESTRICT", "RESTRICT")
		fmt.Printf("OpenUnDb: UserOrg create Foreign Key!\n")
		gdb.Model(&UserOrg{}).AddForeignKey("uid", "user_infos(uid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserOrg{}).AddForeignKey("orid", "organizations(orid)", "RESTRICT", "RESTRICT")
		gdb.Model(&UserOrg{}).AddForeignKey("deid", "departments(deid)", "RESTRICT", "RESTRICT")

		gdb.Model(&RoleObject{}).AddUniqueIndex("idx_roleobject_rid_oid", "rid", "oid")
		gdb.Model(&UserRole{}).AddUniqueIndex("idx_userrole_uid_rid", "uid", "rid")
		gdb.Model(&DepartmentOR{}).AddUniqueIndex("idx_departmentor_deid_oid", "deid", "oid")
		gdb.Model(&UserOrg{}).AddUniqueIndex("idx_userorg_uid_orid_deid", "uid", "orid", "deid")

	}

	fmt.Printf("OnlyOpenUnDbOpenUnDb: db open ok!\n")

	//enforcer.InitCasbin(connect_str)
	return nil
}
