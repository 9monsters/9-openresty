# API 接口参考

## 认证

所有 API（除 `/api/auth/*`）均需 `Authorization: Bearer <token>` 头。

### POST /api/auth/login

**请求体**:
```json
{ "username": "admin", "password": "admin123" }
```

**响应**:
```json
{
  "token": "eyJhbGciOi...",
  "expiresAt": "2026-10-19T02:14:32Z",
  "user": { "id": 1, "username": "admin", "role": "Administrator" }
}
```

### POST /api/auth/refresh

**请求头**: `Authorization: Bearer <existing-token>`

**响应**: 同上格式，签发新 token。

---

## 仪表盘

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/dashboard` | Nginx 状态 + 版本 + 连接统计 |
| GET | `/api/system/info` | CPU/内存/负载/uptime |
| GET | `/api/system/self-check` | 4 项自检 |
| GET | `/api/system/upgrade` | 版本检查 |
| GET | `/api/system/about` | 版本信息 |

## Nginx 控制

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/nginx/status` | running/stopped/unknown |
| POST | `/api/nginx/test` | `nginx -t` |
| POST | `/api/nginx/reload` | 先 test 再 reload |
| POST | `/api/nginx/configure-status` | 部署 stub_status |
| POST | `/api/nginx/restart` | 501 |
| POST | `/api/nginx/stop` | 501 |

## 站点管理

| 方法 | 路径 | 请求/响应 |
|------|------|-----------|
| GET | `/api/sites` | `{sites: [...]}` |
| GET | `/api/sites/:name` | 单个 site 对象 |
| POST | `/api/sites` | `{serverName, listen, root, ssl, locations, luaDirectives}` → `{ok, site, filePath, releaseId}` |
| PUT | `/api/sites/:name` | `{content}` → `{ok, path, releaseId}` |
| DELETE | `/api/sites/:name` | `{ok, name, releaseId}` |

## 配置文件

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/config` | 文件树 |
| GET | `/api/config/read?path=` | 读文件 |
| POST | `/api/config` | `{path, content, overwrite?}` |
| PUT | `/api/config` | `{path, content}` |
| DELETE | `/api/config?path=` | 删除文件 |
| GET | `/api/config/releases` | 发布历史 |
| GET | `/api/config/history?path=` | 文件变更历史 |
| GET | `/api/config/history/:id?path=` | 历史版本内容 |

## 日志

| 方法 | 路径 | 参数 |
|------|------|------|
| GET | `/api/logs/access` | `?lines=200` |
| GET | `/api/logs/error` | `?lines=200` |
| GET | `/api/logs/files` | — |
| GET | `/api/logs/file` | `?path=&lines=200` |

## 通用 CRUD

以下模块均支持标准 `GET / POST / PUT / DELETE`:

| 模块 | 前缀 | 特殊端点 |
|------|------|----------|
| Nodes | `/api/nodes` | |
| Upstreams | `/api/upstreams` | 列表含 `servers` 关联 |
| Streams | `/api/streams` | |
| Certificates | `/api/certificates` | 列表脱敏 certPEM/keyPEM |
| ACME Users | `/api/acme-users` | |
| Snippets | `/api/snippets` | `GET /api/snippets/templates` |
| Access Lists | `/api/access-lists` | |
| Namespaces | `/api/namespaces` | |
| DNS Credentials | `/api/dns-credentials` | secretKey 脱敏 |
| DNS Domains | `/api/dns-domains` | |
| DNS Groups | `/api/dns-groups` | |
| DNS DDNS | `/api/dns-ddns` | |
| Lua Scripts | `/api/lua` | `GET /api/lua/libraries` |

## 用户

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users` | 用户列表 |
| POST | `/api/users` | `{username, email, password, role}` |
| PUT | `/api/users/:id` | `{username, email, role, status, theme, language}` |
| DELETE | `/api/users/:id` | |
| GET | `/api/users/profile` | 当前用户 |
| PUT | `/api/users/profile` | 更新个人资料 |
| PUT | `/api/users/profile/password` | `{currentPassword, newPassword}` |

## 通知

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/notifications` | `?type=&page=1` |
| POST | `/api/notifications` | `{type, title, message}` |
| PUT | `/api/notifications/:id/read` | 标记已读 |
| PUT | `/api/notifications/read-all` | 全部已读 |
| DELETE | `/api/notifications/:id` | |

## 备份

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/backup` | 备份列表 |
| POST | `/api/backup` | `{name?}` |
| POST | `/api/backup/:id/restore` | |
| DELETE | `/api/backup/:id` | |
| GET | `/api/backup/auto-backup` | 自动备份设置 |
| PUT | `/api/backup/auto-backup` | 更新设置 |

## 设置

| 方法 | 路径 |
|------|------|
| GET | `/api/settings` |
| PUT | `/api/settings` |
| GET | `/api/settings/server-name` |
| POST | `/api/settings/configure-status` |

## 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/system` | 系统历史数据（环形缓冲） |
| GET | `/api/stats/nginx` | Nginx 历史数据 |

## 可观测

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/metrics` | Prometheus 指标 |