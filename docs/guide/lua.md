# Lua 脚本

管理 OpenResty Lua 脚本，支持所有 35 种生命周期指令。

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/lua` | Lua 脚本列表 |
| GET | `/api/lua/libraries` | 可用 Lua 库列表 |
| POST | `/api/lua` | 创建脚本 |
| PUT | `/api/lua/:id` | 更新脚本 |
| DELETE | `/api/lua/:id` | 删除脚本 |

## 脚本结构

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

## 生命周期指令

`init_by_lua` → `init_worker_by_lua` → `set_by_lua` → `rewrite_by_lua` →
`access_by_lua` → `content_by_lua` → `header_filter_by_lua` →
`body_filter_by_lua` → `log_by_lua` → ...

## 可用库

| 库 | 说明 |
|----|------|
| `resty.http` | HTTP 客户端 |
| `resty.redis` | Redis 客户端 |
| `resty.mysql` | MySQL 客户端 |
| `resty.websocket` | WebSocket 支持 |
| `resty.limit.traffic` | 流量限制 |
| `resty.core` | 核心增强 |
| `cjson` | JSON 编解码 |