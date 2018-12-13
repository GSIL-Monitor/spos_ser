package enforcer

import (
	"github.com/casbin/casbin"
	"github.com/casbin/gorm-adapter"
)

var (
	cb_en = &casbin.Enforcer{}
)

func NewCasbinEnforcer(connStr string) *casbin.Enforcer {
	Adapter := gormadapter.NewAdapter("mysql", connStr, true)
	enforcer := casbin.NewEnforcer(casbin.NewModel(CasbinConf), Adapter)
	return enforcer
}

func InitCasbin(connect_str string) {
	cb_en = NewCasbinEnforcer(connect_str)
}

func AddPolicy(id, content, action string) bool {
	return cb_en.AddPolicy(id, content, action)
}

func AddRoleForUser(user, role string) bool {
	return cb_en.AddRoleForUser(user, role)
}

func Enforce(uid, path, action string) bool {
	return cb_en.Enforce(uid, path, action)
}
