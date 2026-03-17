# gorbac-starter

- 使用

```go
gocfgmodule.Register(rbac_starter.RbacModule{})
```

- 配置

```ini
[rbac_starter.xxx]
# 设置表名
t_rule= ; sys_auth_rule
t_item= ; sys_auth_item
t_item_child= ; sys_auth_item_child
t_assignment= ; sys_auth_assignment

# 驱动
driver= ; mysql|redis

# 数据源名称（需先初始化 goframework 的 mysql/redis 实例）
table=  ; redis前缀

# redis驱动时的前缀表名
# 是否开启缓存

cache= ; 1|0
# 游客角色名称，默认guest

# 可选：将 RBAC 快照缓存到 Redis（多实例共享）
cache_store= ; redis|空
cache_store_db= ; 默认同 db
cache_store_prefix= ; 默认 gorbac
cache_store_ttl= ; 默认 10m，可写 10m/1h 或秒数

guest=guest
```
```
