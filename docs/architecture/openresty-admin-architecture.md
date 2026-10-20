# OpenResty Admin 架构设计与演进路线

- 状态：建设方案 / 第一阶段实施基线
- 日期：2026-10-16
- 适用范围：当前单机版 OpenResty Admin；多节点能力作为后续演进，不纳入本阶段交付

## 1. 目标与非目标

### 目标

1. 管理操作可追溯、可校验、可回滚，尤其是会影响线上流量的配置变更。
2. 将站点、Lua 和 WAF 放到统一的期望状态与发布模型中，而不是散落的文件编辑入口。
3. 默认限制管理面暴露和高危能力；在认证、授权、审计补齐前维持 loopback 部署边界。
4. 保留现有技术栈和单机易部署性，先通过边界和流程治理提升可靠性。
5. 为未来多节点预留 Agent 边界，但不现在引入分布式系统复杂度。

### 非目标

- 本阶段不提供多租户 SaaS、跨地域集群、自动扩缩容或通用远程 Shell 平台。
- 不承诺任意 Lua 代码可以被安全沙箱化；Lua 发布权限应按高危运维权限处理。
- 不假设所有 OpenResty 构建都安装了某一个特定 WAF 模块或第三方 Lua 库。

## 2. 当前系统基线

当前仓库是单体 Web 管理工具：

- Web：Vue 3、Vite、Pinia、antdv-next、UnoCSS。
- API：Fastify + TypeScript，按路由文件拆分域。
- 状态：SQLite 文件存储；OpenResty 实际配置仍由文件系统持有。
- 运行时操作：服务层读取/写入配置文件、检查 nginx、reload/stop 等。
- Lua：当前主要从 nginx 配置扫描 Lua 指令；库和 shared-dict 页面数据尚非可信运行时探测。
- WAF：仓库尚无完整引擎适配、规则发布和误报处置闭环。

当前关键风险及约束：

- API 尚无可依赖的认证/授权边界；用户密码处理和终端能力也需要安全改造。远程监听必须视为未就绪。
- 配置保存、配置测试和 reload 分散在不同请求中，无法保证测试的内容就是随后生效的内容。
- 站点和通用配置存在直接文件操作入口；必须统一经发布协调服务。
- SQLite 适合单机单实例；当前 collection/JSON payload 结构便于迁移，但不应成为所有领域模型的长期形态。

## 3. 目标逻辑架构

```text
┌────────────────────────────────────────────────────────────┐
│ Vue 管理控制台                                              │
│ 配置草稿 / diff / Lua 包与阶段 / WAF 规则与观察数据         │
└──────────────────────────┬─────────────────────────────────┘
                           │ HTTPS + API schema
┌──────────────────────────▼─────────────────────────────────┐
│ Fastify 控制面                                               │
│ 认证 / RBAC / 参数校验 / API 版本 / 审计上下文              │
├────────────────────────────────────────────────────────────┤
│ 领域模块                                                     │
│ Sites · Certificates · Lua · WAF · DNS · Backups · Users    │
├────────────────────────────────────────────────────────────┤
│ 变更协调器                                                   │
│ Draft → Diff → Validate → Apply → Health → Rollback         │
│ 操作锁 · 任务状态 · 版本记录 · 幂等键                       │
├──────────────────────┬─────────────────────────────────────┤
│ SQLite 控制面状态    │ OpenResty Host Adapter              │
│ 期望对象/版本/审计   │ 文件、证书、nginx -t、reload、指标  │
└──────────────────────┴───────────────────┬─────────────────┘
                                            │ 本机受限操作
                                 ┌──────────▼───────────┐
                                 │ OpenResty 数据面     │
                                 │ NGINX + Lua + WAF    │
                                 └──────────────────────┘
```

### 分层职责

1. **HTTP/API 层**：鉴权、输入 schema、HTTP 状态码、请求身份；不直接拼 Shell、不直接改文件。
2. **领域服务层**：站点、证书、Lua 包、WAF 规则的业务约束和期望状态。
3. **变更协调层**：对变更排序、加锁、记录版本、调用校验/发布/健康检查，并在失败时恢复。
4. **Host Adapter**：唯一允许触碰 OpenResty 文件和进程的边界；路径解析必须防目录穿越和 symlink escape，命令必须参数化执行。
5. **数据层**：SQLite 保存控制面数据和操作记录；运行时配置仍由 OpenResty 文件系统消费。二者通过发布版本关联。

当前可继续采用 Fastify + Vue 3，不建议仅为“换框架”而迁移。演进重点是建立 schema、领域服务、Host Adapter 和统一错误/任务模型。

## 4. 配置发布模型

配置文件仍是 OpenResty 的运行输入，但不能再把“写文件”视为一次完整操作。

### 4.1 发布状态机

```text
Draft
  └─提交发布→ Preparing
       ├─路径/版本冲突→ Rejected
       ├─写入候选失败→ Failed
       └─写入成功→ Validating
            ├─nginx -t 失败→ Restoring → ValidationFailed
            └─测试成功→ Reloading
                 ├─reload 失败→ Restoring → Retesting
                 │                 ├─恢复测试失败→ RecoveryRequired
                 │                 └─恢复测试成功→ ReloadingPrevious → RolledBack
                 └─reload 成功→ HealthChecking
                                      ├─探针失败→ Restoring → RolledBack
                                      └─探针通过→ Succeeded
```

### 4.2 一致性规则

- 一个进程内同一时刻只允许一个配置发布；目前管理面单实例运行，多进程部署时须升级为 SQLite 锁/租约或独立发布 worker。
- 最终目标：每次发布绑定不可变的候选内容摘要和输入版本，避免测试 A、实际发布 B。目前写盘后验证且 reload 前再次读盘比对，不解决跨进程/外部并发写入的所有竞态。
- `nginx -t` 非零退出即失败；绝不能根据输出字符串猜测成功，也不能继续 reload。
- reload 返回成功只代表信号发送成功；后续要检查 master/worker 状态与应用探针。
- 失败时恢复原文件并再次测试；恢复状态不健康时进入 `RecoveryRequired`，明确提示人工处理。当前写盘快照保存在内存，进程崩溃后的自动恢复需后续 WAL/不可变快照支持。
- 删除也属于发布，必须有可恢复快照；目录递归删除不属于配置发布能力。
- 所有操作产生 release/audit 记录，包括操作者、来源、对象、摘要、阶段结果和时间。

### 4.3 版本及漂移

目标数据对象：`config_revision`（不可变内容/摘要）、`release`（一次发布）、`release_item`（文件变化）、`audit_event`（身份与动作）。当前版本可继续让磁盘文件作为权威生效态，逐步加上数据库期望态和磁盘漂移扫描。配置/证书可能含密钥，不应无保护地把完整内容复制进审计表；快照目录需有严格文件权限，未来可加密。

## 5. Lua 脚本与 OpenResty 生命周期

### 5.1 Lua 是版本化制品，不是随意插入的文本

建议 Lua 制品字段：

- `moduleId`、名称、语义化版本、内容摘要、入口文件、依赖及版本。
- 目标 OpenResty/LuaJIT 版本和所需 `lua-resty-*` 库。
- 生命周期函数挂载点、作用范围（全局 / server / location）、启停状态。
- 校验、发布和回滚记录；脚本来源/签名和审批人。

版本不可变；编辑生成新版本。支持 diff、依赖解析、灰度范围和单包回滚。不要在 Web 请求处理中执行用户提交的 Lua 代码。

### 5.2 生命周期挂载模型

| 阶段 | 管理用途 | 约束/风险 |
| --- | --- | --- |
| `init_by_lua*` | 主进程初始化、加载只读配置 | 仅启动/重载阶段运行；避免请求级操作 |
| `init_worker_by_lua*` | worker 初始化、定时任务 | 每个 worker 都运行；定时器要处理重复、退出和重载 |
| `rewrite_by_lua*` | URI/路由调整 | 对请求早期路径产生影响 |
| `access_by_lua*` | 鉴权、限流、WAF | 请求关键路径；限制阻塞、超时和外部依赖 |
| `content_by_lua*` | 动态响应 | 管理连接池、错误响应和请求耗时 |
| `header_filter_by_lua*` | 响应头变更 | 阶段能力受限，不应假设可以 yield |
| `body_filter_by_lua*` | 流式响应体处理 | 需处理分块语义、缓冲和大响应内存 |
| `log_by_lua*` | 请求后记录 | 保持低开销；不要把同步外部调用放在收尾路径 |
| `balancer_by_lua*` | 动态选择 upstream peer | 可用 API 和上下文受限，单独测试 |
| SSL 阶段 | 动态证书或 TLS 行为 | 影响握手；单独授权、探针和回滚 |

发布流水线至少包含 Lua 语法/静态检查、目标版本兼容性、依赖检测、临时 OpenResty 实例验证、`nginx -t`、灰度/健康检查和回滚。Linter 规则需配置 OpenResty globals。不能把静态检查当沙箱或安全证明。

### 5.3 与配置的关系

Lua 包自身保存代码和版本；部署适配器生成受控的 `*_by_lua_file` 挂载或 include 文件，并写入对应 NGINX 上下文。生命周期阶段、上下文合法性、作用范围由 schema 限定，不能允许模块声明任意 nginx 指令片段绕过校验。模块更新也走配置发布状态机。

## 6. WAF 管理

WAF 分为引擎、规则集、站点绑定、排除项、事件/指标五个模型：

- **Engine adapter**：探测构建和模块能力；对 Lua WAF、ModSecurity/CRS 等分别做适配，不假设模块默认存在。
- **Rule set**：名称、来源、版本、签名/摘要、规则数、兼容范围、启用状态。
- **Binding**：规则集绑定的站点/路径/阶段及优先级。
- **Exception**：必须有规则 ID、受限站点/路径/参数范围、原因、创建者、过期时间；禁止默认全局永久白名单。
- **Event**：规则 ID/版本、请求关联 ID、站点、检测动作、最终处置、时间；敏感请求数据按最小化原则记录。

运行策略采用 `off → observe → block`：新规则先观察，只记日志不阻断；确认误报率可接受后按站点或流量范围灰度阻断。发布需展示冲突检查、影响范围、告警阈值及回滚方式。核心指标包括拦截率、误报/排除率、Top 规则、规则执行耗时和加载错误。

## 7. 安全、权限和运维

### 安全基线（远程监听前置）

1. 所有 API 默认认证；密码使用成熟 KDF 哈希，令牌有过期和撤销策略。
2. RBAC 将只读、配置编辑、发布、证书私钥、WAF 管理、用户管理、终端拆成独立能力。
3. 高危发布要求重新认证/二次确认；审计记录身份和不可抵赖上下文。
4. Terminal 默认关闭；如保留需独立开关、强授权、短期会话和完整审计。不要让通用 API 进程以无边界 root 权限运行。
5. Secrets 使用 OS keychain/外部密钥服务；SQLite、日志、导出和备份不得意外泄漏。
6. 网络保持 loopback，直到上面能力完成并经过安全测试。

### 主机运行

- Host Adapter 使用明确的可执行文件和参数数组，限制路径白名单、超时、并发和操作类型。
- 将服务账号能力降到最低；避免 Fastify API 直接获得任意 sudo。
- 统一结构化任务状态：queued/running/succeeded/failed/rolled_back/recovery_required。
- 监控管理 API、数据库、磁盘空间、OpenResty 状态、证书有效期、配置漂移和发布失败。
- SQLite WAL 仅用于单机单写者场景；做好原子备份和恢复演练，不把 DB 文件放在网络共享盘。

## 8. 数据模型与 API 约定

### 核心实体（目标）

- `User`, `Role`, `Permission`, `Session`。
- `Node`（当前单机也使用本地 node 记录，为后续 Agent 预留）。
- `ConfigDraft`, `ConfigRevision`, `Release`, `ReleaseItem`, `AuditEvent`。
- `LuaPackage`, `LuaPackageVersion`, `LuaMount`。
- `WafEngine`, `WafRuleSet`, `WafBinding`, `WafException`, `WafEvent`。
- `OperationTask`, `SecretReference`, `BackupSnapshot`。

SQLite 中重要领域实体应逐步从泛型 JSON payload 转成有约束的表；用迁移脚本演进 schema，不在所有功能完成前先引入大型 ORM。

### API 原则

- OpenAPI/schema-first：请求、响应、错误结构和权限声明明确。
- 读操作返回期望态、实际态及版本；写操作返回 operation/release ID。
- 耗时操作支持 `GET /operations/:id` 查询，不依赖长 HTTP 请求保持连接。
- 变更支持幂等键和版本冲突检测；错误使用稳定错误码。

## 9. 框架建议

- **继续 Fastify**：插件封装和 schema 能力适合该规模，先建立 route → domain service → host adapter 结构。
- **继续 Vue 3 / Vite / Pinia**：以 API schema/生成类型作为前后端契约来源；UI 将“保存草稿”和“发布”明确分开。
- **继续 SQLite**：适用于本地安装；加版本迁移、备份校验和最小权限。明确不支持多 API 实例并发写。
- **配置渲染**：先使用领域对象 + 明确模板/序列化器，避免在多个 route 手写字符串拼接。生成结果必须可预览、可 diff、可回滚。
- **测试**：纯逻辑单测、Host Adapter 假实现、临时 OpenResty 集成测试、API 权限测试和发布故障注入分层建设。

## 10. 分阶段路线图

| 阶段 | 交付 | 验收重点 |
| --- | --- | --- |
| P0 安全与配置安全 | loopback、身份/RBAC 设计、Host Adapter、配置统一发布、审计 | 测试失败绝不 reload；故障恢复可验证；权限测试通过 |
| P1 运营可靠性 | release 历史、diff、健康探针、漂移检测、数据库迁移/备份恢复 | 从历史版本恢复演练成功 |
| P2 Lua 管理 | 版本化包、生命周期挂载、依赖/兼容检查、临时实例验证 | 阶段不合法会被拒绝；失败发布自动回滚 |
| P3 WAF | 引擎探测、规则集版本、observe、例外到期、灰度 block | 误报可定位，范围和回滚清楚 |
| P4 多节点（有真实需求再做） | 本地 Agent、mTLS、期望状态、批次发布 | Agent 受限权限、离线重试、节点级漂移 |

本次第一阶段实现边界：提供单进程配置发布协调器和 SQLite 摘要发布记录；通用配置、站点和 status 配置写入均经过“快照 → 写入 → `nginx -t` → reload”，失败恢复和重新校验。默认不自动备份完整历史内容（只有单份 `.bak`，包含密钥时权限必须受控）；配置列表中的历史页面仍需要单独做多版本存储。`/nginx/reload` 是显式手动操作，先执行 `nginx -t`，不具备“与发布协调器同锁”的保证；`/nginx/stop` 和 `/nginx/restart` 暂返回 501，以避免无认证状态下绕过发布闭环。此阶段不代表身份/RBAC、跨进程锁、健康探针、Lua/WAF runtime 或多节点 Agent 已完成。

## 11. 关键决策记录

### ADR-001：先采用模块化单体

**决策**：保持 Fastify/Vue/SQLite 单机部署；以领域服务和 Host Adapter 建立边界。

**原因**：当前产品定位是单机管理工具，拆微服务会放大部署、身份和一致性复杂度，不能解决当前发布安全问题。

### ADR-002：所有配置变更共享发布协调器

**决策**：站点、通用配置以及 Lua/WAF 生成配置统一使用一条发布状态机。

**原因**：多个入口独立写文件会产生不同的备份、校验和回滚语义；应只允许一个 Host Adapter 管理文件和 reload。

### ADR-003：Lua 与 WAF 作为版本化领域对象

**决策**：Lua 包、生命周期挂载、WAF 规则集和例外具备独立 schema，但其运行态变更仍进入统一发布流程。

**原因**：它们有各自的兼容、灰度和审计需求，但不能绕开配置校验和主机权限边界。

## 12. 第一阶段实现与操作约束

- `services/config-release.ts` 是发布入口；通过依赖注入的 Host Adapter 做单机发布/回滚测试，SQLite `config-releases` 记录来源、文件路径和内容哈希，不记录文件原文或操作者（认证未实现，不能伪造身份）。
- `routes/config.ts`、`routes/sites.ts` 和 status 配置使用同一协调器；新建站点遇到文件已存在时拒绝覆盖。HTTP `422` 为校验失败，`409` 为发布忙/文件冲突，`500` 为写入/回滚问题；成功包含 `releaseId`。`GET /api/config/releases` 返回最近 50 条元数据，不包含密钥正文。
- 执行环境：`getNginxPaths()` 的当前路径；`nginx -t -c <nginxConfPath>`，非零退出不 reload；reload 优先 `sudo -n ... -s reload -c <nginxConfPath>`，不可用时尝试无 sudo，同样依据退出码。主机命令通过 argv 调用，不用 shell 拼接。
- `writeConfigFile` 的 `.bak` 是最近一次旧内容，不能等同多版本备份；删除文件的恢复由协调器用预写内存快照实现。symlink 和目录都不是本阶段允许修改的配置文件。
- 目前脚本运行、Lua 阶段挂载和 WAF 规则只是架构设计，不要把文档中的目标态写成已实现功能。使用真实 OpenResty 主机前必须先解决认证和凭据问题；测试只在临时目录使用注入 Host Adapter，不触碰宿主机服务。
- 验证：`pnpm check`、`pnpm --filter @openresty-admin/server test`；真实主机变更必须遵守仓库 AGENTS.md 的保护顺序与人工授权。
