# 访问控制

## 功能概览

管理 IP 访问控制列表（ACL），支持 allow/deny 规则。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/access-lists` | ACL 列表 |
| POST | `/api/access-lists` | 创建 ACL |
| PUT | `/api/access-lists/:id` | 更新 ACL |
| DELETE | `/api/access-lists/:id` | 删除 ACL |

ACL 行结构：

```json
{
  "id": 1,
  "name": "office-whitelist",
  "description": "Office IP whitelist",
  "rules": [
    { "id": 1, "type": "allow", "ip": "10.0.0.0", "cidr": 8 },
    { "id": 2, "type": "deny", "ip": "10.0.1.100" }
  ]
}
```

# Lua 脚本

## 功能概览

管理 OpenResty Lua 脚本，支持所有 35 种生命周期指令。

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/lua` | Lua 脚本列表 |
| GET | `/api/lua/libraries` | 可用 Lua 库列表 |
| POST | `/api/lua` | 创建脚本 |
| PUT | `/api/lua/:id` | 更新脚本 |
| DELETE | `/api/lua/:id` | 删除脚本 |

脚本行结构：

```json
{
  "id": 1,
  "name": "rate-limiter",
  "directive": "access_by_lua_block",
  "value": "local limit = ngx.shared.limit\n...",
  "filePath": "/etc/nginx/lua/rate.lua",
  "siteName": "example.com"
}
```

库列表：resty.http, resty.redis, resty.mysql, resty.websocket, resty.limit.traffic, resty.core, cjson（均标为 installed）。