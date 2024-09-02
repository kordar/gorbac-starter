package gorbac_starter

import (
	goframeworkgoredis "github.com/kordar/goframework-goredis"
	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	logger "github.com/kordar/gologger"
	"github.com/kordar/gorbac"
	"github.com/kordar/gorbac-gorm"
	"github.com/kordar/gorbac-redis"
	"github.com/spf13/cast"
)

var (
	authManager gorbac.AuthManager
)

func GetAuthManager() gorbac.AuthManager {
	return authManager
}

func getMapStr(m map[string]interface{}, field string, value string) string {
	if m[field] == "" {
		return value
	} else {
		return cast.ToString(m[field])
	}
}

type RbacModule struct {
}

func (m RbacModule) Name() string {
	return "rbac_starter"
}

func (m RbacModule) Load(value interface{}) {
	cfg := cast.ToStringMap(value)
	gorbac.SetTableName("rule", getMapStr(cfg, "t_rule", "sys_auth_rule"))
	gorbac.SetTableName("item", getMapStr(cfg, "t_item", "sys_auth_item"))
	gorbac.SetTableName("item-child", getMapStr(cfg, "t_item_child", "sys_auth_item_child"))
	gorbac.SetTableName("assignment", getMapStr(cfg, "t_assignment", "sys_auth_assignment"))

	driver := getMapStr(cfg, "driver", "mysql")

	var repos gorbac.AuthRepository
	db := getMapStr(cfg, "db", "gorbac")
	if driver == "mysql" {
		if !goframeworkgormmysql.HasMysqlInstance(db) {
			logger.Warnf("初始化rbac组件失败，请先初始化数据库%s", db)
			return
		}

		mysqlDB := goframeworkgormmysql.GetMysqlDB(db)
		repos = gorbac_gorm.NewSqlRbac(mysqlDB)
	}

	if driver == "redis" {
		if !goframeworkgoredis.HasRedisInstance(db) {
			logger.Warnf("初始化rbac组件失败，请先初始化数据库%s", db)
			return
		}
		tb := getMapStr(cfg, "table", "gorbac_table")
		mysqlDB := goframeworkgoredis.GetRedisClient(db)
		repos = gorbac_redis.NewRedisRbac(mysqlDB, tb)
	}

	cache := cast.ToBool(cfg["cache"])
	rbacManager := gorbac.NewDefaultManager(repos, cache)

	guest := getMapStr(cfg, "guest", "guest")
	role := rbacManager.CreateRole(guest)
	rbacManager.SetDefaultRoles(role)
	authManager = rbacManager
}

func (m RbacModule) Close() {
}
