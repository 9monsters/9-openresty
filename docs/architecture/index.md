# 架构概览

## 整体架构

OpenResty Admin 是**单机** OpenResty/NGINX 管理工具：一个 Go 二进制内嵌 Vue 3 前端，通过 REST API 管理本机 OpenResty 实例。

```
┌──────────────────────────────────────────────────┐
│           Browser (Vue 3 SPA, 嵌入)               │
└────────────────────┬─────────────────────────────┘
                     │ HTTP /api/* (JWT + RBAC)
┌────────────────────▼─────────────────────────────┐
│            openresty-admin (Go 二进制)             │
│                                                    │
│  Gin Router ── 认证/三角色 RBAC 中间件              │
│      │                                             │
│      ├─ handler/*      —— 各域 CRUD                 │
│      ├─ release/*      —— 发布协调器（状态机）       │
│      ├─ nginx/*        —— 唯一 Host Adapter          │
│      ├─ stats/*        —— 内存环形缓冲区             │
│      └─ model/*        —— GORM 模型                 │
│                    │                               │
│      ┌─────────────▼──────────────┐                │
│      │ PostgreSQL 或 SQLite (GORM) │                │
│      └────────────────────────────┘                │
└──────────────────────────────────────────────────┘
                     │ 文件 + 进程（受限操作）
┌────────────────────▼─────────────────────────────┐
│              本机 OpenResty / NGINX                │
│      配置文件 · nginx -t · reload · 日志 · 指标      │
└──────────────────────────────────────────────────┘
```

## 技术栈

| 层 | 技术 | 说明 |
|----|------|------|
| HTTP 框架 | **Gin** | 路由、中间件、JSON 绑定 |
| ORM | **GORM** | AutoMigrate、类型化查询 |
| 数据库 | **PostgreSQL 14+ 或 SQLite** | 双驱动，`DATABASE_DRIVER` 切换 |
| CLI | **Cobra** | 命令行参数解析 |
| 认证 | **JWT (HMAC-SHA256) + bcrypt** | 密码哈希 + 三角色 RBAC |
| 配置发布 | **Release Coordinator** | 快照→写入→nginx -t→reload→健康探针→回滚 |
| 日志 | **slog** | 结构化日志 |
| 指标 | **Prometheus client** | /metrics 端点 |
| 前端 | **Vue 3 + Vite + antdv-next** | 嵌入二进制，或独立部署 |
| WAF | **Lua-based** | `access_by_lua*` 注入 |

## 数据流

```
请求流:
  Browser → Gin Router → Auth → RBAC → Handler → GORM → DB

配置发布流 (经发布协调器):
  Browser → POST /api/config → Coordinator.Publish()
    → 快照 → 写入 → nginx -t → reload → HTTP 健康探针
    → 失败自动回滚 / 成功记录 release

统计流:
  Collector(1s) → RingBuffer[600] → GET /api/stats/* → ECharts 图表
```

## ADR 列表

| # | 决策 | 日期 | 状态 |
|---|------|------|------|
| [003](/architecture/adr-003-jwt-auth) | Web 端 JWT 认证 | 2026-10-17 | 现行 |
| [004](/architecture/adr-004-postgresql) | PostgreSQL 主用 | 2026-10-17 | 现行（已扩展 SQLite） |
| [007](/architecture/adr-007-observability) | slog + trace + Prometheus | 2026-10-17 | 现行 |
| [008](/architecture/adr-008-roadmap) | 实施路线 | 2026-10-17 | 现行 |
| [009](/architecture/adr-009-schema) | 规范化数据库 Schema | 2026-10-18 | 现行 |

多机集群相关 ADR（001/002/005/006）已归档到 `architecture/archived/`。
