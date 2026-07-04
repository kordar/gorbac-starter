package gorbacstarterfx

import (
	"fmt"
	"log/slog"
	"time"

	gocfgmodulefx "github.com/kordar/gocfg-load-module/fx/v2"
	goframeworkgoredis "github.com/kordar/goframework-goredis"
	goframeworkgormmysql "github.com/kordar/goframework-gorm-mysql"
	"github.com/kordar/gorbac"
	gorbac_cache_redis "github.com/kordar/gorbac-cache-redis"
	gorbac_gorm "github.com/kordar/gorbac-gorm"
	gorbac_redis "github.com/kordar/gorbac-redis"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

// ModuleConfig gorbac 模块配置
type ModuleConfig struct {
	DB              string        // MySQL/Redis 实例名称, 默认 "gorbac"
	Driver          string        // "mysql" 或 "redis", 默认 "mysql"
	GuestRole       string        // 游客角色名称, 默认 "guest"
	CacheEnabled    bool          // 是否开启内存缓存
	CacheStore      string        // "redis" 启用 Redis 分布式缓存
	CacheStoreDB    string        // 缓存 Redis 实例名称, 默认同 DB
	CachePrefix     string        // 缓存 key 前缀, 默认 "gorbac"
	CacheTTL        time.Duration // 缓存过期时间, 默认 10 分钟
	TableRule       string        // 规则表名, 默认 "sys_auth_rule"
	TableItem       string        // 权限项表名, 默认 "sys_auth_item"
	TableItemChild  string        // 权限项关联表名, 默认 "sys_auth_item_child"
	TableAssignment string        // 权限分配表名, 默认 "sys_auth_assignment"
	RedisTable      string        // Redis 驱动时表前缀, 默认 "gorbac_table"
}

func normalizeConfig(c ModuleConfig) ModuleConfig {
	if c.DB == "" {
		c.DB = "gorbac"
	}
	if c.Driver == "" {
		c.Driver = "mysql"
	}
	if c.GuestRole == "" {
		c.GuestRole = "guest"
	}
	if c.CachePrefix == "" {
		c.CachePrefix = "gorbac"
	}
	if c.CacheTTL <= 0 {
		c.CacheTTL = 10 * time.Minute
	}
	if c.CacheStoreDB == "" {
		c.CacheStoreDB = c.DB
	}
	if c.TableRule == "" {
		c.TableRule = "sys_auth_rule"
	}
	if c.TableItem == "" {
		c.TableItem = "sys_auth_item"
	}
	if c.TableItemChild == "" {
		c.TableItemChild = "sys_auth_item_child"
	}
	if c.TableAssignment == "" {
		c.TableAssignment = "sys_auth_assignment"
	}
	if c.RedisTable == "" {
		c.RedisTable = "gorbac_table"
	}
	return c
}

var _ gocfgmodulefx.GoCfgModule = cfgModule{}
var _ gocfgmodulefx.GoCfgIndex = cfgModule{}

type cfgModule struct {
	name  string
	index int
}

type Option func(*cfgModule)

// WithIndex 设置优先级
func WithIndex(index int) Option {
	return func(s *cfgModule) {
		s.index = index
	}
}

// StarterModule 返回可注册到 gocfg-load-module/fx 的 gorbac 模块适配器。
func StarterModule(name string, opts ...Option) gocfgmodulefx.GoCfgModule {
	c := &cfgModule{name: name}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (m cfgModule) Name() string {
	return m.name
}

func (m cfgModule) Index() int {
	return m.index
}

func (m cfgModule) Load(data any) []fx.Option {
	slog.Info("Module Load Complete", "module", "gorbac-starter(fx)")
	return []fx.Option{
		Module(buildModuleConfig(data)),
	}
}

// buildModuleConfig 从 gocfg section 数据构建 ModuleConfig。
// 支持 [gorbac] 段内的 flat key-value 配置。
func buildModuleConfig(data any) ModuleConfig {
	cfg := ModuleConfig{}
	m := cast.ToStringMapString(data)
	if len(m) == 0 {
		return cfg
	}

	cfg.Driver = m["driver"]
	cfg.DB = m["db"]
	cfg.GuestRole = m["guest_role"]
	cfg.CacheStore = m["cache_store"]
	cfg.CacheStoreDB = m["cache_store_db"]
	cfg.CachePrefix = m["cache_prefix"]
	cfg.TableRule = m["table_rule"]
	cfg.TableItem = m["table_item"]
	cfg.TableItemChild = m["table_item_child"]
	cfg.TableAssignment = m["table_assignment"]
	cfg.RedisTable = m["redis_table"]

	if v := m["cache_enabled"]; v == "true" || v == "1" {
		cfg.CacheEnabled = true
	}
	if v := m["cache_ttl"]; v != "" {
		if d, err := cast.ToDurationE(v); err == nil {
			cfg.CacheTTL = d
		}
	}

	return cfg
}

// Module 返回一个 fx.Option, 按配置初始化并注册 *gorbac.RbacService 到 fx 容器。
func Module(config ModuleConfig) fx.Option {
	config = normalizeConfig(config)

	return fx.Module("gorbac-starter",
		fx.Supply(config),
		fx.Provide(provideRbacService),
	)
}

func provideRbacService(config ModuleConfig) (*gorbac.RbacService, error) {
	gorbac.SetTableName("rule", config.TableRule)
	gorbac.SetTableName("item", config.TableItem)
	gorbac.SetTableName("item-child", config.TableItemChild)
	gorbac.SetTableName("assignment", config.TableAssignment)

	var repos gorbac.AuthRepository

	switch config.Driver {
	case "mysql":
		if !goframeworkgormmysql.HasMysqlInstance(config.DB) {
			return nil, fmt.Errorf("gorbac: mysql instance %q not found", config.DB)
		}
		mysqlDB := goframeworkgormmysql.GetMysqlDB(config.DB)
		repos = gorbac_gorm.NewSqlRbac(mysqlDB)
	case "redis":
		if !goframeworkgoredis.HasRedisInstance(config.DB) {
			return nil, fmt.Errorf("gorbac: redis instance %q not found", config.DB)
		}
		redisDB := goframeworkgoredis.GetRedisClient(config.DB)
		repos = gorbac_redis.NewRedisRbac(redisDB, config.RedisTable)
	default:
		return nil, fmt.Errorf("gorbac: unsupported driver %q", config.Driver)
	}

	rbacManager := gorbac.NewDefaultManager(repos, config.CacheEnabled)

	if config.CacheEnabled && config.CacheStore == "redis" {
		if !goframeworkgoredis.HasRedisInstance(config.CacheStoreDB) {
			return nil, fmt.Errorf("gorbac: cache redis instance %q not found", config.CacheStoreDB)
		}
		storeRedis := goframeworkgoredis.GetRedisClient(config.CacheStoreDB)
		store := gorbac_cache_redis.NewRedisCacheStore(storeRedis)
		rbacManager.SetCacheStore(store, config.CachePrefix, config.CacheTTL)
	}

	role := rbacManager.CreateRole(config.GuestRole)
	rbacManager.SetDefaultRoles(role)

	rbacService := gorbac.NewRbacServiceWithManager(rbacManager)
	slog.Info("gorbac initialized", "driver", config.Driver, "db", config.DB)
	return rbacService, nil
}
