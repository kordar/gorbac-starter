package gorbac_starter

import (
	"log/slog"
	"time"

	goframeworkgoredis "github.com/kordar/goframework-goredis"
	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	"github.com/kordar/gorbac"
	gorbac_cache_redis "github.com/kordar/gorbac-cache-redis"
	gorbac_gorm "github.com/kordar/gorbac-gorm"
	gorbac_redis "github.com/kordar/gorbac-redis"
	"github.com/spf13/cast"
)

var (
	rbacservice *gorbac.RbacService
)

func GetRbacService() *gorbac.RbacService {
	return rbacservice
}

func GetAuthManager() gorbac.AuthManager {
	return rbacservice.GetAuthManager()
}

func getMapStr(m map[string]interface{}, field string, value string) string {
	if m[field] == nil {
		return value
	} else {
		return cast.ToString(m[field])
	}
}

func getDuration(m map[string]interface{}, field string, value time.Duration) time.Duration {
	if m[field] == nil {
		return value
	}
	raw := m[field]
	if s, ok := raw.(string); ok {
		if s == "" {
			return value
		}
		if d, err := time.ParseDuration(s); err == nil {
			return d
		}
	}
	sec := cast.ToInt64(raw)
	if sec <= 0 {
		return value
	}
	return time.Duration(sec) * time.Second
}

type RbacModule struct {
	name    string
	load    func(name string, value map[string]interface{})
	depends []string
}

func NewRbacModule(name string, load func(moduleName string, item map[string]interface{}), depends ...string) *RbacModule {
	return &RbacModule{name, load, depends}
}

func (m RbacModule) Depends() []string {
	return m.depends
}

func (m RbacModule) Name() string {
	return m.name
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
			slog.Warn("初始化rbac组件失败，请先初始化数据库", "module", m.Name(), "db", db)
			return
		}

		mysqlDB := goframeworkgormmysql.GetMysqlDB(db)
		repos = gorbac_gorm.NewSqlRbac(mysqlDB)
	}

	if driver == "redis" {
		if !goframeworkgoredis.HasRedisInstance(db) {
			slog.Warn("初始化rbac组件失败，请先初始化数据库", "module", m.Name(), "db", db)
			return
		}
		tb := getMapStr(cfg, "table", "gorbac_table")
		redisDb := goframeworkgoredis.GetRedisClient(db)
		repos = gorbac_redis.NewRedisRbac(redisDb, tb)
	}

	if repos == nil {
		slog.Warn("初始化rbac组件失败，未识别的driver", "module", m.Name(), "driver", driver)
		return
	}

	cache := cast.ToBool(cfg["cache"])
	rbacManager := gorbac.NewDefaultManager(repos, cache)

	if cache && getMapStr(cfg, "cache_store", "") == "redis" {
		cacheDB := getMapStr(cfg, "cache_store_db", db)
		if !goframeworkgoredis.HasRedisInstance(cacheDB) {
			slog.Warn("初始化rbac缓存失败，请先初始化redis", "module", m.Name(), "redis", cacheDB)
			return
		}
		prefix := getMapStr(cfg, "cache_store_prefix", "gorbac")
		ttl := getDuration(cfg, "cache_store_ttl", 10*time.Minute)
		storeRedis := goframeworkgoredis.GetRedisClient(cacheDB)
		store := gorbac_cache_redis.NewRedisCacheStore(storeRedis)
		rbacManager.SetCacheStore(store, prefix, ttl)
	}

	guest := getMapStr(cfg, "guest", "guest")
	role := rbacManager.CreateRole(guest)
	rbacManager.SetDefaultRoles(role)
	//
	rbacservice = gorbac.NewRbacServiceWithManager(rbacManager)

	if m.load != nil {
		m.load(m.Name(), cfg)
	}
}

func (m RbacModule) Close() {
}
