# Stream / Upstream

## Stream（TCP/UDP 四层代理）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/streams` | 列出所有 stream |
| POST | `/api/streams` | 创建 stream |
| PUT | `/api/streams/:id` | 更新 stream |
| DELETE | `/api/streams/:id` | 删除 stream |

Stream 行结构：

```json
{
  "id": 1,
  "name": "mysql-proxy",
  "listen": "3306",
  "protocol": "tcp",
  "proxyPass": "192.168.1.10:3306",
  "ssl": false,
  "enabled": true
}
```

## Upstream（七层负载均衡）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/upstreams` | 列出 upstream 组（含服务器列表） |
| POST | `/api/upstreams` | 创建 upstream 组 |
| PUT | `/api/upstreams/:id` | 更新 upstream 组 |

Upstream 行结构：

```json
{
  "id": 1,
  "name": "backend-api",
  "strategy": "round_robin",
  "servers": [
    { "id": 1, "address": "192.168.1.10", "port": 8080, "weight": 1, "maxFails": 3, "failTimeout": "30s", "status": "up" }
  ]
}
```

支持的策略：`round_robin`、`least_conn`、`ip_hash`。