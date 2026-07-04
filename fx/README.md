# gorbac-starter/fx

基于 uber-go/fx 的 RBAC 启动器，DI 注入 `*gorbac.RbacService`，支持 **MySQL / Redis** 双驱动及 Redis 分布式缓存。

## 安装

```bash
go get github.com/kordar/gorbac-starter/fx/v2
```

## 配置

```ini
[gorbac]
driver          = mysql       ; mysql | redis，默认 mysql
db              = sys         ; 预初始化的 goframework 实例名，默认 "gorbac"
guest_role      = guest       ; 游客角色，默认 "guest"
cache_enabled   = true        ; 开启内存缓存
cache_store     = redis       ; Redis 分布式缓存
cache_store_db  = sys         ; 缓存 Redis 实例名，默认同 db
cache_prefix    = gorbac      ; 缓存 key 前缀，默认 "gorbac"
cache_ttl       = 10m         ; 缓存过期，默认 10m
table_rule      = sys_auth_rule
table_item      = sys_auth_item
table_item_child = sys_auth_item_child
table_assignment = sys_auth_assignment
redis_table     = gorbac_table ; Redis 驱动 key 前缀
```

## 注册到 Fx

```go
import (
    gorbacstarterfx "github.com/kordar/gorbac-starter/fx/v2"
    gocfgmodulefx "github.com/kordar/gocfg-load-module/fx/v2"
    "github.com/kordar/gorbac"
    "go.uber.org/fx"
)

func AdminServerStarter() *fx.App {
    return ServerWithStarter(
        []gocfgmodulefx.GoCfgModule{
            gorbacstarterfx.StarterModule("gorbac", gorbacstarterfx.WithIndex(2)),
        },
        nil, nil,
        fx.Invoke(func(svc *gorbac.RbacService) {
            // 自动注入，无需标签
        }),
    )
}
```

## 驱动依赖

| 驱动 | 前置模块 | 底层 |
|------|---------|------|
| `mysql` | mysql-starter | [gorbac-gorm](https://github.com/kordar/gorbac-gorm) |
| `redis` | goredis-starter | [gorbac-redis](https://github.com/kordar/gorbac-redis) |

## 配置项参考

| Key | 类型 | 默认 | 说明 |
|-----|------|------|------|
| `driver` | string | `mysql` | 存储驱动 |
| `db` | string | `gorbac` | 数据源实例名 |
| `guest_role` | string | `guest` | 游客角色 |
| `cache_enabled` | bool | `false` | 内存缓存 |
| `cache_store` | string | — | `redis` 启用分布式缓存 |
| `cache_store_db` | string | 同 `db` | 缓存 Redis 实例 |
| `cache_prefix` | string | `gorbac` | 缓存 key 前缀 |
| `cache_ttl` | duration | `10m` | 缓存过期 |
| `table_rule` | string | `sys_auth_rule` | 规则表 |
| `table_item` | string | `sys_auth_item` | 权限项表 |
| `table_item_child` | string | `sys_auth_item_child` | 权限关联表 |
| `table_assignment` | string | `sys_auth_assignment` | 权限分配表 |
| `redis_table` | string | `gorbac_table` | Redis key 前缀 |

## Option

```go
gorbacstarterfx.StarterModule("gorbac",
    gorbacstarterfx.WithIndex(2), // 优先级
)
```

## 模块接口

```go
// GoCfgModule
func Name() string
func Load(data any) []fx.Option

// GoCfgIndex（可选）
func Index() int
```

## License

MIT
