# OpenResty Admin — 功能模块清单

## 1. 认证 (Auth)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/auth/login` | POST | 用户名密码登录，返回 JWT |
| `/api/auth/refresh` | POST | 刷新即将过期的 token |

**后端**: `internal/auth/jwt.go` (HMAC-SHA256 签名), `internal/auth/middleware.go` (HTTP 中间件), `internal/api/auth.go` (login/refresh handler)
**前端**: `views/login/LoginPage.vue` (登录卡片), `api/auth.ts`

---

## 2. 仪表盘 (Dashboard)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/dashboard` | GET | Nginx 状态 + 版本 + 连接统计 + 服务器时间 |

**后端**: `internal/api/system.go` → `handleDashboard`
**前端**: `views/dashboard/` — HomePage, ServerDash, NginxDash, SiteDash

---

## 3. 系统信息 (System)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/system/info` | GET | CPU/内存/负载/平台/uptime |
| `/api/system/self-check` | GET | Nginx 二进制 + 配置语法 + 权限检查 |
| `/api/system/upgrade` | GET | 版本检查（当前 0.1.0，桩） |
| `/api/system/about` | GET | 版本信息 (Go/Nginx/平台/数据库类型) |

**后端**: `internal/api/system.go` + `internal/api/sysinfo.go` (内存/负载/uptime 采集)
**前端**: `views/system/` — HomePage, SelfCheck, Upgrade, About

---

## 4. Nginx 控制 (Nginx Process)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/nginx/status` | GET | running/stopped/unknown + 版本 |
| `/api/nginx/test` | POST | `nginx -t` 配置测试 |
| `/api/nginx/reload` | POST | 先 `nginx -t`，通过后 reload |
| `/api/nginx/configure-status` | POST | 通过 release 协调器写入 stub_status 配置 |
| `/api/nginx/restart` | POST | HTTP 501（暂不支持） |
| `/api/nginx/stop` | POST | HTTP 501（暂不支持） |

**后端**: `internal/nginx/host.go` — 二进制操作、状态检测、统计采集

---

## 5. 站点管理 (Sites)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/sites` | GET | 列出所有 server block |
| `/api/sites/{name}` | GET | 按 server_name 查单个站点 |
| `/api/sites` | POST | 新建 server block，通过 release 协调器写入 |
| `/api/sites/{name}` | PUT | 更新站点配置内容 |
| `/api/sites/{name}` | DELETE | 删除站点配置文件 |

**后端**: `internal/api/sites.go` — CRUD + nginx config 渲染 (`buildSiteConfig`)
**前端**: `views/sites/SiteList.vue`, `SiteEdit.vue`
**关键**: 所有写入走 `release.Coordinator.Publish()` — snapshot → write → test → reload，失败自动回滚

---

## 6. 配置文件 (Config Files)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/config` | GET | 列出 nginx 配置目录树 |
| `/api/config/read` | GET | `?path=` 读取文件内容 |
| `/api/config` | POST | 创建新配置文件 |
| `/api/config` | PUT | 更新配置文件 |
| `/api/config` | DELETE | `?path=` 删除配置文件 |
| `/api/config/releases` | GET | 发布历史记录 |
| `/api/config/history` | GET | `?path=` 查看文件备份历史 |
| `/api/config/history/{id}` | GET | 查看某个备份的内容 |

**后端**: `internal/api/config.go` + `internal/release/coordinator.go` (发布协调器) + `internal/nginx/backup.go` (备份管理)
**前端**: `views/config/ConfigEditor.vue`

---

## 7. 日志 (Logs)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/logs/access` | GET | 最后 N 行 access.log |
| `/api/logs/error` | GET | 最后 N 行 error.log |
| `/api/logs/files` | GET | 列出日志目录下的 .log 文件 |
| `/api/logs/file` | GET | `?path=` 读取指定日志文件 |

**后端**: `internal/api/logs.go`
**前端**: `views/logs/` — LogViewer (access/error 实时查看), LogList

---

## 8. Lua 管理

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/lua` | GET | 列出所有 config 中的 Lua 代码块 |
| `/api/lua/libraries` | GET | OpenResty Lua 库列表（桩数据） |
| `/api/lua/shared-dicts` | GET | lua_shared_dict 声明（桩） |

**后端**: `internal/api/lua.go` + `internal/nginx/parser.go` (Lua directive 提取)
**前端**: `views/lua/LuaManager.vue`

---

## 9. 终端 (Terminal)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/terminal/ws` | WebSocket | 双向 shell 代理，cwd = nginx conf 目录 |

**后端**: `internal/api/terminal.go` — gorilla/websocket + exec.Command("/bin/bash")
**前端**: `views/terminal/Terminal.vue` — xterm.js

---

## 10. 用户管理 (Users)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/users` | GET | 用户列表（脱敏） |
| `/api/users` | POST | 创建用户 |
| `/api/users/{id}` | PUT | 更新用户 |
| `/api/users/{id}` | DELETE | 删除用户 |
| `/api/users/profile` | GET | 当前用户资料 |
| `/api/users/profile` | PUT | 更新个人资料 |
| `/api/users/profile/password` | PUT | 修改密码 |

**后端**: `internal/api/users.go`
**前端**: `views/user/UserList.vue`, `UserProfile.vue`

---

## 11. 通知 (Notifications)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/notifications` | GET | 通知列表，支持 `?type=` 筛选 + 分页 |
| `/api/notifications` | POST | 创建通知 |
| `/api/notifications/{id}/read` | PUT | 标记已读 |
| `/api/notifications/read-all` | PUT | 全部已读 |
| `/api/notifications/{id}` | DELETE | 删除通知 |

**后端**: `internal/api/notifications.go`
**前端**: `views/notification/NotificationList.vue`

---

## 12. 备份 (Backup)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/backup` | GET | 备份列表 |
| `/api/backup` | POST | 创建手动备份 |
| `/api/backup/{id}/restore` | POST | 恢复备份（桩） |
| `/api/backup/{id}` | DELETE | 删除备份 |
| `/api/backup/auto-backup` | GET | 自动备份设置 |
| `/api/backup/auto-backup` | PUT | 更新自动备份设置 |

**后端**: `internal/api/backup.go`
**前端**: `views/backup/` — BackupHome, AutoBackup

---

## 13. 设置 (Settings)

| 路由 | 方法 | 说明 |
|------|------|------|
| `/api/settings` | GET | 所有设置 + nginx 路径 + 状态 |
| `/api/settings` | PUT | 更新设置并实时应用 (`UpdatePaths`) |
| `/api/settings/server-name` | GET | 仅服务器名称 |
| `/api/settings/configure-status` | POST | 部署 stub_status 配置 |

**后端**: `internal/api/settings.go`
**前端**: `views/settings/SettingsPage.vue`, `NginxHostSetup.vue`

---

## 14. CRUD 集合（通用模式）

以下模块共享 `internal/api/collections.go` 中的 `registerCRUD` 泛型路由注册器，
每个模块自动获得 `GET/POST/PUT/DELETE` 标准 CRUD：

| 集合 | 路由前缀 | DB 表名 | 数据项 |
|------|----------|---------|--------|
| Upstream | `/api/upstreams` | `upstreams` | 策略 + 服务器列表 |
| Stream | `/api/streams` | `streams` | TCP/UDP 代理 |
| 证书 | `/api/certificates` | `certificates` | PEM/域名/到期时间（列表页脱敏） |
| 命名空间 | `/api/namespaces` | `namespaces` | 逻辑隔离 |
| 节点 | `/api/nodes` | `nodes` | Agent 注册节点（状态 + 端口 + 环境） |
| ACME 用户 | `/api/acme-users` | `acme-users` | Let's Encrypt 账户 |
| 代码片段 | `/api/snippets` | `snippets` | 可复用配置片段（列表截断至 100 字符） |
| 访问控制 | `/api/access-lists` | `access-lists` | IP 规则列表 |
| DNS 凭证 | `/api/dns-credentials` | `dns-credentials` | Cloudflare/阿里云 API key（脱敏） |
| DNS 域名 | `/api/dns-domains` | `dns-domains` | 域名区域 |
| DNS 分组 | `/api/dns-groups` | `dns-groups` | 路由策略 |
| DDNS | `/api/dns-ddns` | `dns-ddns` | 动态 DNS 端点 |

**后端**: `internal/api/collections.go` + `internal/api/seeds.go` (种子数据)
**前端**: 各自 `api/*.ts` + `views/*/`

---

## 15. 集群与 Agent（gRPC）

| Service | RPC | 说明 |
|---------|-----|------|
| `AgentRegistration` | `Register`, `Heartbeat` | Agent 注册 + 心跳 |
| `ConfigDeploy` | `Push`, `Validate`, `Reload`, `Rollback` | 配置下发 |
| `StatsCollector` | `StreamStats` | Agent → Hub 指标流 |
| `TerminalProxy` | `Shell` | 双向终端代理 |
| `LogReader` | `ReadLog` | 远程日志读取 |

**后端**: `proto/openresty/admin/v1/services.proto` → `gen/go/...`
**Hub 端**: `internal/hub/manager.go` (连接池), `internal/hub/registration.go`, `internal/hub/server.go`, `internal/hub/metrics.go`
**Agent 端**: `cmd/openresty-agent/main.go`, `internal/agent/agent.go`, `internal/agent/services.go`
**部署协调**: `internal/release/cluster.go` — 集群扇出 `ClusterCoordinator`

---

## 16. 基础设施

| 包 | 文件 | 职责 |
|----|------|------|
| `internal/config` | `config.go` | 配置加载：CLI > 文件 > 环境变量 |
| `internal/store` | `store.go` | PostgreSQL JSONB CRUD，版本化迁移，JSON 导入 |
| `internal/nginx` | `host.go` | Nginx Host Adapter — 二进制操作 + 配置 I/O |
| `internal/nginx` | `parser.go` | Server/Location block 解析 + Lua 指令提取 |
| `internal/nginx` | `backup.go` | 配置备份快照 |
| `internal/release` | `coordinator.go` | 单机发布协调器 (snapshot → write → test → reload) |
| `internal/release` | `cluster.go` | 集群发布协调器 (并行扇出 + 聚合) |
| `internal/auth` | `jwt.go`, `middleware.go`, `key.go` | JWT 签发/验证/中间件 |
| `internal/trace` | `trace.go` | UUIDv7 trace ID + HTTP 中间件 |
| `internal/metrics` | `metrics.go` | Prometheus 指标 (HTTP/gRPC/deploy) |
| `internal/adapters` | `host.go`, `remote.go` | NginxHost 接口 + RemoteHost gRPC 代理桩 |
| `internal/embedweb` | `embed.go` | 嵌入 Vue 前端 + SPA fallback |