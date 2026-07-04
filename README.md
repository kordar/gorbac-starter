# gorbac-starter

基于 [gorbac](https://github.com/kordar/gorbac) 的 RBAC 启动器，支持 **MySQL / Redis 双驱动**、缓存加速、GoCfg 模块化加载。

本仓库包含两个包：

| 包 | 路径 | 适用场景 |
|---|------|---------|
| 基础包 | `github.com/kordar/gorbac-starter` | 通用模块，对接 gocfg-load-module |
| Fx 包 | `github.com/kordar/gorbac-starter/fx/v2` | 对接 uber-go/fx，DI 注入 `*gorbac.RbacService` |

## 安装

```bash
go get github.com/kordar/gorbac-starter
```

## 基础包用法

### 定义模块

```go
import (
    gorbacstarter "github.com/kordar/gorbac-starter"
    gocfgmodule "github.com/kordar/gocfg-load-module"
)

mod := gorbacstarter.NewRbacModule(
    "gorbac",
    func(moduleName string, item map[string]interface{}) {
        // 加载完成回调
        log.Printf("%s loaded", moduleName)
    },
    "mysql", // 依赖模块（需先初始化 mysql/redis）
)
```

### 注册并加载

```go
gocfgmodule.Register(mod, "mysql") // gorbac 依赖 mysql
gocfgmodule.RefreshDepends(nil)
gocfgmodule.ResolveAll(settings)
```

### 配置

```ini
[gorbac]
driver      = mysql       ; mysql | redis
db          = sys         ; 预初始化的实例名
guest       = guest       ; 游客角色名

; 表名（可选，默认值如下）
t_rule       = sys_auth_rule
t_item       = sys_auth_item
t_item_child = sys_auth_item_child
t_assignment = sys_auth_assignment

; 缓存（可选）
cache       = true             ; 开启内存缓存
cache_store = redis            ; Redis 分布式缓存
cache_store_db = sys           ; 缓存 Redis 实例
cache_store_prefix = gorbac    ; key 前缀
cache_store_ttl = 10m          ; 过期时间

; Redis 驱动专用
table       = gorbac_table     ; key 前缀
```

### 获取已初始化的实例

```go
svc := gorbacstarter.GetRbacService()
mgr := gorbacstarter.GetAuthManager()
```

## 驱动支持

| 驱动 | 底层 |
|------|------|
| `mysql` | [gorbac-gorm](https://github.com/kordar/gorbac-gorm) — MySQL 持久化 |
| `redis` | [gorbac-redis](https://github.com/kordar/gorbac-redis) — Redis 持久化 |

### 缓存加速

| 缓存方式 | 说明 |
|---------|------|
| 内存缓存 (`cache=true`) | 单实例，进程内 |
| Redis 缓存 (`cache_store=redis`) | 多实例共享，需先初始化 goredis 实例 |

## 模块接口

```go
// GoCfgModule
func Name() string
func Load(value interface{})
func Close()

// Depends（可选）
func Depends() []string
```

## License

MIT
