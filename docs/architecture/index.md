# 架构概览

## 整体架构

```
┌─────────────────────────────────┐
│         Browser (Vue 3 SPA)     │
└──────────────┬──────────────────┘
               │ HTTP (JSON API)
┌──────────────▼──────────────────┐
│     Hub (openresty-admin)       │
│  ┌────────┐  ┌────────────────┐ │
│  │  Gin   │  │  gRPC Server   │ │
│  │ Router │  │    (:9090)     │ │
│  └───┬────┘  └───────┬────────┘ │
│      │               │          │
│  ┌───▼───────────────▼────────┐ │
│  │        Service Layer       │ │
│  │  auth / handler / nginx    │ │
│  └───────────┬────────────────┘ │
│              │                  │
│  ┌───────────▼────────────────┐ │
│  │    PostgreSQL (GORM)       │ │
│  └────────────────────────────┘ │
└─────────────────────────────────┘
               │ gRPC
     ┌─────────┼─────────┐
     ▼         ▼         ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Agent 1 │ │ Agent 2 │ │ Agent N │
│ (go)    │ │ (go)    │ │ (go)    │
│ nginx   │ │ nginx   │ │ nginx   │
└─────────┘ └─────────┘ └─────────┘
```

## 技术栈

| 层 | 技术 | 说明 |
|----|------|------|
| HTTP 框架 | **Gin** | 路由、中间件、JSON 绑定 |
| ORM | **GORM** | AutoMigrate、类型化查询 |
| 数据库 | **PostgreSQL 14+** | 30+ 张表，JSONB 用于 config_releases |
| CLI | **Cobra** | 命令行参数解析 |
| 认证 | **JWT (HMAC-SHA256)** | 自实现，无第三方库依赖 |
| gRPC | **protobuf + grpc-go** | Hub-Agent 通信 |
| 日志 | **slog** | Go 1.21+ 标准库结构化日志 |
| 指标 | **Prometheus client** | /metrics 端点 |
| 前端 | **Vue 3 + Vite + antdv-next** | SPA，嵌入 Go 二进制 |
| 终端 | **gorilla/websocket + xterm.js** | WebSocket shell |

## 数据流

```
请求流:
  Browser → Gin Router → Auth Middleware → Handler → GORM → PostgreSQL

配置变更流:
  Browser → POST /api/sites → Handler → Build nginx config
    → host.WriteConfig() 写入文件
    → host.Test()    nginx -t
    → host.Reload()   nginx -s reload
    → GORM.Create(config_history)

统计流:
  Collector(1s) → RingBuffer[600] → GET /api/stats/system → ECharts 图表

Agent 流:
  Agent.Start() → gRPC Register → Heartbeat(15s)
  Hub → gRPC PushConfig → Agent.WriteConfig → nginx -t → reload → Report
```

## ADR 列表

| # | 决策 | 日期 |
|---|------|------|
| [001](/architecture/adr-001-cluster) | 多机集群 Hub-Agent 架构 | 2026-10-17 |
| [002](/architecture/adr-002-grpc) | gRPC 通信 + mTLS | 2026-10-17 |
| [003](/architecture/adr-003-jwt-auth) | Web 端 JWT 认证 | 2026-10-17 |
| [004](/architecture/adr-004-postgresql) | PostgreSQL 替代 SQLite | 2026-10-17 |
| [005](/architecture/adr-005-fat-agent) | 胖 Agent 模型 | 2026-10-17 |
| [006](/architecture/adr-006-incremental-deploy) | 增量配置部署 | 2026-10-17 |
| [007](/architecture/adr-007-observability) | slog + trace + Prometheus | 2026-10-17 |
| [008](/architecture/adr-008-roadmap) | 自底向上 5 阶段路线 | 2026-10-17 |
| [009](/architecture/adr-009-schema) | 规范化数据库 Schema | 2026-10-18 |