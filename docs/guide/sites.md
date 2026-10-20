# 站点管理

## 功能概览

管理 Nginx server block 的全生命周期：创建、编辑、删除、启用/禁用。

## 路由

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/sites` | 列出所有站点（解析配置文件） |
| GET | `/api/sites/:name` | 按 server_name 查询 |
| POST | `/api/sites` | 创建新 server block |
| PUT | `/api/sites/:name` | 更新站点配置 |
| DELETE | `/api/sites/:name` | 删除站点配置文件 |

## 创建站点

POST `/api/sites` 请求体：

```json
{
  "serverName": ["example.com", "www.example.com"],
  "listen": ["80", "443 ssl http2"],
  "root": "/var/www/html",
  "ssl": true,
  "locations": [
    { "path": "/api/", "proxyPass": "http://backend:8080" },
    { "path": "/static/", "root": "/var/www/static" }
  ],
  "luaDirectives": [
    { "type": "access_by_lua_block", "value": "ngx.say('Hello')" }
  ]
}
```

服务端自动渲染为：

```nginx
server {
  listen 80;
  listen 443 ssl http2;
  server_name example.com www.example.com;
  ssl on;
  ssl_certificate certs/example.com.pem;
  ssl_certificate_key certs/example.com.key;
  root /var/www/html;

  location /api/ {
    proxy_pass http://backend:8080;
  }

  location /static/ {
    root /var/www/static;
  }

  access_by_lua_block {
    ngx.say('Hello')
  }
}
```

## 安全机制

所有写入操作通过 Release Coordinator：
1. 快照原文件 → 2. 写入新内容 → 3. `nginx -t` 验证 → 4. reload
失败自动回滚。所有变更记录在 `config_releases` 表中。