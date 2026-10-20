# 备份与恢复

## 手动备份

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/backup` | 备份列表 |
| POST | `/api/backup` | 创建手动备份 |
| POST | `/api/backup/:id/restore` | 恢复备份 |
| DELETE | `/api/backup/:id` | 删除备份 |

## 自动备份

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/backup/auto-backup` | 自动备份设置 |
| PUT | `/api/backup/auto-backup` | 更新自动备份设置 |

自动备份设置：

```json
{
  "enabled": true,
  "interval": "daily",
  "retention": 7,
  "time": "03:00",
  "path": "/var/backups/openresty"
}
```

## 配置历史

每次通过 Web UI 修改配置文件时自动记录历史版本：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/config/history?path=` | 文件变更历史 |
| GET | `/api/config/history/:id?path=` | 查看历史版本内容 |

# 系统设置

## 设置 API

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/settings` | 所有设置 |
| PUT | `/api/settings` | 更新设置 |
| GET | `/api/settings/server-name` | 仅服务器名称 |

返回结构：

```json
{
  "serverName": "OpenResty Admin",
  "language": "zh_CN",
  "theme": "auto",
  "nginxBin": "/usr/local/bin/openresty",
  "nginxConfDir": "/etc/nginx",
  "paths": { "nginxBin": "...", ... },
  "nginxStatus": "running",
  "version": "openresty/1.31.2.1",
  "isOpenResty": true
}
```

## 运行时 Nginx 路径更新

修改 settings 中的 nginx 路径字段后，Host Adapter 立即切换操作目标，**无需重启服务**。

# 用户与认证

## 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/users` | 用户列表 |
| POST | `/api/users` | 创建用户 |
| PUT | `/api/users/:id` | 更新用户 |
| DELETE | `/api/users/:id` | 删除用户 |
| GET | `/api/users/profile` | 个人资料 |
| PUT | `/api/users/profile` | 更新个人资料 |
| PUT | `/api/users/profile/password` | 修改密码 |

角色：`Administrator`、`Operator`、`Viewer`

## JWT 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/auth/login` | 登录，返回 JWT |
| POST | `/api/auth/refresh` | 刷新 Token |

Token 格式：HMAC-SHA256，24h 有效期。所有 `/api/*` 路径（除 `/api/auth/*`）需 `Authorization: Bearer <token>` 头。