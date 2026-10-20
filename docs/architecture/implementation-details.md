# 模块实现细节

## 1. 认证 (Auth)

### 后端

**JWT 签发** (`internal/auth/jwt.go:116` 行)
- 算法：HMAC-SHA256，标准 3 段式 `header.payload.signature`
- Claims 结构：`uid`、`un`(username)、`role`、`exp`(24h)、`iat`、`jti`(UUIDv7)
- 签名密钥：`auth.SetKey()` 设置 32 字节，默认从 `JWT_SIGNING_KEY` 环境变量读取，若无则自动生成随机密钥（ephemeral，重启后所有 token 失效）
- Token 格式：`eyJhbGci...eyJ1aWQiOjEs...MWLy-73TPlMsmDf-4...`

**中间件** (`internal/auth/middleware.go:67` 行)
- `auth.Middleware()` 包装 `http.Handler`
- 跳过 `/api/auth/*`（登录不拦截自己）、`/metrics`、`/health`、非 API 路径（前端静态文件）
- 从 `Authorization: Bearer <token>` 提取 token → `auth.Verify()` 验签 + 检查过期
- 提取的 `*Claims` 注入 `r.Context()`，通过 `auth.ClaimsFromContext(ctx)` 读取
- 失败返回 401 + `WWW-Authenticate: Bearer realm="OpenResty Admin"`

**Login/Refresh API** (`internal/api/auth.go:82` 行)
- `POST /api/auth/login` — 接受 `{username, password}`，查 `users` 集合明文比对密码，返回 `{token, expiresAt, user}`

**组合链** (main.go:103)：`auth.Middleware(trace.HTTPMiddleware(mux))`，trace 先于 auth 执行

### 前端

**登录页** (`views/login/LoginPage.vue` 162 行)
- 独立页面，不含 DefaultLayout shell — 路由 `/login` meta 标记 `public: true`
- 用户名预填 `admin`，密码空
- 调用 `api/auth.ts` 的 `login()` → `userStore.login(token, user)` 保存到 localStorage
- 语言切换：settings store 的 zh-CN/en-US

**路由守卫** (`routes/index.ts`)
- `beforeEach` 读取 `localStorage.getItem('openresty-token')`
- 无 token + 非 public 路由 → 跳转 `/login`，带 `?redirect=` 参数
- 已有 token + 访问 `/login` → 跳转 `/`
- 401 拦截器 (`api/client.ts`) 清除 token 并跳 `#/login`

### 待实现
- 密码 bcrypt 哈希（当前明文比对）
- Token 自动续期（`/api/auth/refresh` 已有端点，前端未调用）
- 精细化 RBAC 权限检查（当前只验证 token 有效性，未校验 role）


## 2. 数据库存储 (Store)

### 架构 (`internal/store/store.go` 354 行)

**PgSQL + JSONB 文档存储**
- 连接池：`pgxpool`，DSN 来自 `DATABASE_URL` env 或默认 `postgres://openresty:...@192.168.5.100:5432/openresty`
- 版本化迁移：`schema_migrations` 表追踪已执行版本，3 个迁移版本：
  1. `collections` 表 (collection TEXT, seq INT, id BIGINT, data JSONB, PK)
  2. `idx_collections_id` 索引
  3. `meta` 表 (key TEXT PK, value TEXT)

**CRUD API**（全部带 `context.Context` 参数）
| 方法 | SQL | 事务 |
|------|-----|------|
| `Load(collection)` | `SELECT data FROM collections WHERE collection=$1 ORDER BY seq` | 无 |
| `Save(collection, items)` | `DELETE` + 批量 `INSERT` | 有 |
| `Create(collection, item)` | Load → maxID → append → Save | Save 内带事务 |
| `Get(collection, id)` | Load → 遍历 | 无 |
| `Update(collection, id, patch)` | Mutate → Load → 遍历修改 → Save | Save 内带事务 |
| `Delete(collection, id)` | Load → 过滤 → Save | Save 内带事务 |
| `Mutate(collection, fn)` | Load → fn 回调 → Save | Save 内带事务 |

**JSON 导入**：启动时扫描 `data/` 目录下的 `*.json` 文件，按 `meta` 表中标记决定是否导入（最多一次）

**ID 提取** (`NumberID`)：兼容 `json.Number`、`float64`、`int64`、`int` 四种 shape，确保 JS 到 Go 的数值不丢精度

### 当前数据库
- Host: `192.168.5.100:5432`，用户 `openresty`，库 `openresty`
- 初始化脚本：`scripts/init-db.sql`
- 测试 DSN：`OPENRESTY_PGTEST_DSN` 环境变量


## 3. Nginx Host Adapter

### 核心文件 (`internal/nginx/host.go` 466 行)

**Nginx 进程操作**
- `Version()` — 执行 `openresty -v`，判断是否 OpenResty
- `Status()` — 先读 PID 文件 → `kill(pid, 0)`，再 fallback `pgrep nginx: master`
- `Test()` — `nginx -t -c <confPath>`，30s 超时
- `Reload()` — 优先 `sudo -n nginx -s reload`，fallback 直接 reload
- `Stats()` — 先 `curl http://127.0.0.1/nginx_status` 解析 stub_status，失败返回 nil

**配置 I/O**
- `resolveConfigPath()` — 防护链条：拒绝绝对路径 → 拒绝 `..` → 拒绝符号链接 → 拒绝目录
- `ReadConfig()` / `ReadForRelease()` — 通过 `resolveConfigPath` 校验后 `os.ReadFile`
- `WriteConfig()` — 先 snapshot 到 `<name>.bak`，再写入
- `DeleteConfig()` — 拒绝目录，仅删单文件
- `ListConfigFiles()` — `filepath.WalkDir` 遍历，跳过 `.` 开头 / `~` / `.bak`

**路径管理**
- `Paths()` / `UpdatePaths()` — 读写锁保护，支持运行时通过 settings API 修改
- 默认值：brew 安装路径 (`/usr/local/opt/openresty/...`)

### 解析器 (`internal/nginx/parser.go` 207 行)
- `ParseSites(content, filePath)` → `[]ParsedSite`
- 使用有限状态机遍历大括号嵌套，逐字符定位 `server { ... }` 和 `location { ... }`
- 提取 server_name、listen、ssl、http2、root、proxy_pass、access_log、error_log
- `ExtractLuaDirectives()` — 35 种 OpenResty Lua 指令 (by_lua_block / by_lua_file)

### 备份 (`internal/nginx/backup.go` 91 行)
- `BackupConfig()` — 写入 `<confDir>/.backups/<id>_<safeName>` + `<...>.meta`
- `ListBackups()` — 列出指定文件的所有备份，按时间倒序


## 4. 发布协调器 (Release Coordinator)

### 单机版 (`internal/release/coordinator.go` 322 行)

**状态机**：running → persisted → validated → reloaded → succeeded / rolled_back / failed / recovery_required

**Publish 流程**：
1. `mu.TryLock()` — 阻止并发发布
2. `validate(changes)` — 拒绝空路径 / 路径逃逸 / 重复 / 删除不存在文件
3. 快照 `originals` — 调用 `host.ReadForRelease()` 读取每个文件的当前内容
4. `createRecord()` — 写入 `config-releases` 集合 (status=running, beforeHash/afterHash)
5. `apply()` — 逐个写入 changes，已变更的记入 `changed` 列表；失败立即回滚
6. `host.Test()` — `nginx -t`，失败回滚
7. 二次快照检查 — 防止测试期间被外部编辑器篡改
8. `host.Reload()` — reload，失败回滚
9. `finish()` — 更新 release 记录为 succeeded

**回滚**：逆序恢复 `changed` → `nginx -t` 验证恢复后的配置 → 若 reloadPrevious=true 则 reload → 记录最终状态

**测试覆盖** (`coordinator_test.go` 241 行)：8 个测试覆盖成功/验证失败/reload 失败/写磁盘失败/不可恢复/删除回滚/noop/外部篡改

### 集群版 (`internal/release/cluster.go` 123 行)
- `ClusterCoordinator.Publish()` — 并行扇出到所有 agent，每节点一个 goroutine
- 聚合结果：任意节点失败 → `StatusFailed`，全部成功 → `StatusSucceeded`
- 记录格式：`{releaseID, source, nodes[], files, node_results[{nodeId, ok, output, error}]}`


## 5. 站点管理 (Sites)

### 路由 (`internal/api/sites.go` 187 行)

**List** — 遍历所有配置文件 → 逐文件 `ParseSites()` → 合并返回
**Get** — 按 server_name 搜索，匹配第一个
**Create** — 用户提交参数 → `buildSiteConfig()` 渲染完整 server block → `release.Publish()` 写入
**Update** — 按 server_name 找到文件 → release.Publish 覆盖
**Delete** — backup → release.Publish(Delete:true)

### 配置渲染 (`buildSiteConfig`)
- 输入：`{serverName[], listen[], root?, ssl?, luaDirectives?, locations[]}`
- 输出：完整的 nginx server block 字符串
- SSL：自动生成 `ssl_certificate certs/<name>.pem` / `ssl_certificate_key`
- Location：支持 proxy_pass / root / luaDirectives 三层嵌套
- Lua 指令：block 格式用 `{}` 包裹，file 格式用 `;` 结尾

### 前端
- `SiteList.vue` (84 行) — 表格展示，支持搜索和状态筛选
- `SiteEdit.vue` (212 行) — 表单：域名、端口、SSL、root、locations、Lua 指令


## 6. 配置文件管理 (Config)

### 路由 (`internal/api/config.go` 142 行)

- List/Read：直接调 `host.ListConfigFiles/ReadConfig`
- Create/Update：backup 原文件 → release.Publish 写入
- Delete：backup → release.Publish(Delete:true)
- Releases：`release.List()` 返回最近 50 条记录
- History：`host.ListBackups(path)` + 内容查看

### 前端：`ConfigEditor.vue` (211 行)
- 左侧文件树 + 右侧代码编辑区（Monaco 或原生 textarea）
- 发布流程：Save → nginx -t → 显示结果 → Reload


## 7. 仪表盘 (Dashboard)

### API (`internal/api/system.go` 125 行)
- `GET /api/dashboard` → 聚合 `host.Status()` + `host.Version()` + `host.Stats()` + `serverTime`

### 前端 (4 个视图)
- **HomePage** (108 行) — 四卡片入口：Server Status、Nginx Performance、Sites、Logs
- **ServerDash** (78 行) — CPU/内存/负载/uptime，数据来自 `/api/system/info`
- **NginxDash** — 连接统计：Active/Accepts/Handled/Requests/Reading/Writing/Waiting，来自 `/api/dashboard`
- **SiteDash** (75 行) — 站点导航卡片，按 SSL 状态分组


## 8. 系统信息 (System)

### 后端 (`internal/api/system.go` + `internal/api/sysinfo.go` 159 行)

**跨平台指标采集**
| 指标 | macOS | Linux |
|------|-------|-------|
| 内存 | `sysctl hw.memsize` + `vm_stat` 计算 free pages | `/proc/meminfo` (MemTotal/MemAvailable) |
| CPU 型号 | `sysctl machdep.cpu.brand_string` | `/proc/cpuinfo` model name |
| 负载 | `sysctl vm.loadavg` | `/proc/loadavg` |
| Uptime | `sysctl kern.boottime` 计算差值 | `/proc/uptime` |

**Self Check** — 4 项检查：
1. Nginx Binary — `exec.LookPath()` 查找可执行文件
2. Config Syntax — `nginx -t`
3. Write Permission — `os.Stat()` config 目录
4. Management API — running on `host:port`

### 前端
- `HomePage.vue` — 仪表盘入口
- `SelfCheck.vue` — 4 项检查结果卡片（ok/warning/error）
- `Upgrade.vue` — 版本对比（当前为 0.1.0 桩）
- `About.vue` — 版本信息面板


## 9. 日志查看 (Logs)

### 后端 (`internal/api/logs.go` 42 行)
- `ReadLog(kind)` → `host.ReadLog(kind, lines)`，默认 200 行
- `LogFiles()` → `host.LogFiles()`，返回 .log 文件列表
- `ReadLogByPath()` → 直接读取指定日志文件

### 前端
- `LogViewer.vue` (91 行) — 实时 access/error 日志查看，可配行数
- `LogList.vue` — 浏览所有日志文件


## 10. 终端 (Terminal)

### 后端 (`internal/api/terminal.go` 75 行)
- WebSocket 升级：`gorilla/websocket.Upgrader{CheckOrigin: func() bool { return true }}`
- 启动 `/bin/bash`，`cwd = nginxConfDir`
- 双向 pipe：`stdout → ws.WriteMessage()` / `stderr → ws.WriteMessage()` / `ws.ReadMessage() → stdin`
- `io.Copy` 两个 goroutine 并发传输

### 前端：`Terminal.vue` (292KB bundle)
- xterm.js + fit addon
- WebSocket 连接 `/api/terminal/ws`


## 11. 用户管理 (Users)

### 后端 (`internal/api/users.go` 188 行)

**种子数据**：`admin / admin123 / Administrator`（首次调用自动创建）

**路由**：
- `GET /api/users` — 列表，密码字段已脱敏
- `POST /api/users` — 创建，必填 username + email + password
- `PUT /api/users/{id}` — 更新（password 不可通过此端点修改）
- `DELETE /api/users/{id}` — 删除
- `GET /api/users/profile` — 当前用户（取第一个用户）
- `PUT /api/users/profile` — 更新个人资料
- `PUT /api/users/profile/password` — 修改密码（验证旧密码）

### 前端
- `UserList.vue` (125 行) — 表格 CRUD，角色下拉（Administrator/Operator/Viewer）
- `Profile.vue` (97 行) — 表单：邮箱/语言 + 修改密码


## 12. CRUD 集合

### 后端 (`internal/api/collections.go` 293 行)

**泛型路由注册器** `registerCRUD(mux, crudSpec)`
- 每个集合自动获得 5 个路由：GET / GET by ID / POST / PUT / DELETE
- 配置项：`collection`(表名)、`listKey`(JSON 键)、`itemKey`、`detail`(是否支持单条查询)、`seed`(初始数据)、`prepare`(创建预处理)、`decorate`(读取后处理)

**prepare 函数示例**
- `prepareUpstream`：默认 strategy=round_robin，给每个 server 填充 weight/maxFails/failTimeout
- `prepareCertificate`：默认 issuer=Manual，autoRenew=false，gen PEM/key/notBefore/notAfter
- `prepareNode`：默认 port=9001，status=offline，environment=production
- `prepareAcme`：默认 provider=Let's Encrypt，status=active
- `prepareSnippet`：填充 createdAt
- `prepareAccessList`：默认 rules=[]
- `prepareDNSGroup`：默认 domains=[]

**decorate 函数**（读取脱敏）
- `stripCertificateSecrets`：删除 certPEM / keyPEM
- `stripSecretKey`：删除 secretKey
- `stripPassword`：删除 password
- `truncateSnippet`：content 截断至 100 字符

### 前端
每个集合对应 `api/<domain>.ts` + `views/<domain>/<List>.vue` 模式
- antdv-next `ATable` + `APopconfirm` 删除确认 + `message.success/error` 反馈
- `onMounted` + `loading` ref + `ASpin` 加载态


## 13. 通知 (Notifications)

### 后端 (`internal/api/notifications.go` 131 行)
- 种子：4 条示例（success 加载/high 证书即将过期/error upstream down/info 备份完成）
- 分页：`?page=1&pageSize=20`，`?type=success|warning|error|info`
- 批量已读：`PUT /api/notifications/read-all` 遍历设置 `read=true`

### 前端 (`NotificationList.vue` 88 行)
- Tab 切换（All/Success/Warning/Error/Info）
- 未读标记 + 单条/批量标记已读


## 14. 备份 (Backup)

### 后端 (`internal/api/backup.go` 136 行)
- 种子：daily-2026-10-07(auto) + manual-2026-10-06(manual)
- Create：生成随机 size(100KB-600KB) 和 files(5-25)
- Restore：stub，返回 `Restored from <name>`
- Auto-backup settings：`{enabled, interval, retention, time, path}`，默认 daily/7d/03:00

### 前端 (`views/backup/`)
- `BackupHome.vue` (140 行) — 备份列表 + 创建/删除
- `AutoBackup.vue` (62 行) — 设置表单


## 15. 设置 (Settings)

### 后端 (`internal/api/settings.go` 88 行)
- 默认值合并：`defaultSettings()` → 从 DB 加载覆盖 → 返回
- `UpdatePaths` 实时生效：settings API 修改 nginx 路径后立即通过 `host.UpdatePaths()` 更新 Host Adapter
- 保存在 `app-settings` 集合

### 前端 (`views/settings/`)
- `SettingsPage.vue` (132 行) — 主题/语言/Nginx 路径
- `NginxHostSetup.vue` — SSH 连接配置


## 16. 集群与 Agent

### Proto 定义 (`proto/openresty/admin/v1/services.proto`)
- 5 个 gRPC service，14 个 RPC
- 生成：`scripts/gen-proto.sh` → `gen/go/openresty/admin/v1/`

### Hub 端
- `internal/hub/manager.go` — 连接池 (`map[string]*Connection`)，RWMutex 保护，Register/Remove/Get/List/Count/FanOut
- `internal/hub/registration.go` — 注册 handler
- `internal/hub/server.go` — gRPC 服务器启动/关闭
- `internal/hub/metrics.go` — 每 10s 更新 `agentsOnline` gauge

### Agent 端
- `cmd/openresty-agent/main.go` — cobra CLI，`--hub-addr`/`--token`/`--metrics-port`
- `internal/agent/agent.go` — 启动 metrics HTTP server + 心跳循环
- `internal/agent/services.go` — `ConfigDeployService`(push/validate/rollback journal) + `StatsService` + `TerminalService` + `LogService`

### 增量部署流程
1. Hub 创建 release record
2. 转换 `nginx.Change` → `adapters.AgentChange`
3. 并行 gRPC `PushConfig` 到每个 agent
4. Agent 写文件 → `nginx -t` → reload → 报告结果
5. Hub 聚合结果，任一失败则回滚已成功的节点


## 17. 可观测性

### Trace (`internal/trace/trace.go` 31 行)
- UUIDv7 trace ID，HTTP middleware 注入/提取 `X-Trace-Id` header
- 注入到 `context.Context`，贯穿整个请求生命周期

### Metrics (`internal/metrics/metrics.go`)
- 7 个 Prometheus 指标：
  - `http_requests_total` (counter, by method/path/status)
  - `http_request_duration_seconds` (histogram)
  - `agents_online` (gauge)
  - `deploy_total` (counter, by status)
  - `deploy_duration_seconds` (histogram)
  - `agent_rpc_duration_seconds` (histogram, by agent/method)
  - `agent_rpc_errors_total` (counter, by agent/method)
- `/metrics` 端点通过 `promhttp.Handler` 暴露

### 日志
- `log/slog` JSON handler，stdout 输出
- 每行含 `time`/`level`/`msg` + 上下文属性 + `trace_id`


## 18. 前端架构

**入口链：** `main.ts` → `App.vue`(AConfigProvider + theme) → `router` → `DefaultLayout.vue` → `RouterView`

**路由** (`routes/modules/all.ts`)
- 21 个 route module，`path` + `name` + `component` + `meta: {name, icon?, hiddenInSidebar?, hideChildren?}`
- 侧栏菜单：DefaultLayout 从路由自动生成 `menuItems`，icon 通过 `iconMap` 映射

**状态管理** (Pinia)
- `user` store — token/info/login/logout
- `settings` store — theme/language/nginxStatus/serverName

**i18n** — 简单字典，zh-CN 有 ~230 条翻译，en-US 回退到 key 本身

**设计语言** — iOS 26 Liquid Glass：`ios-glass`/`ios-group`/`ios-row`/`ios-instrument`/`ios-chip`/`ios-mono`