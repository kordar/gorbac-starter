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
table=  ; redis前缀
# 是否开启缓存
cache= ; 1|0
# 游客角色名称，默认guest
guest=guest
```