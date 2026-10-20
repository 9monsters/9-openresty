# 部署概览

## 拓扑

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Browser    │────▶│  Hub (:3000) │────▶│ PostgreSQL   │
│   (Vue SPA)  │     │  + gRPC 9090 │     │ 192.168.5.100│
└──────────────┘     └──────┬───────┘     └──────────────┘
                            │ gRPC mTLS
              ┌─────────────┼─────────────┐
              ▼             ▼             ▼
        ┌──────────┐ ┌──────────┐ ┌──────────┐
        │ Agent 1  │ │ Agent 2  │ │ Agent N  │
        │ :9101    │ │ :9101    │ │ :9101    │
        │ nginx    │ │ nginx    │ │ nginx    │
        └──────────┘ └──────────┘ └──────────┘
              ▲             ▲             ▲
              └─────────────┼─────────────┘
                     Prometheus scrape
```

## 端口规划

| 端口 | 服务 | 协议 |
|------|------|------|
| 3000 | Hub HTTP API + Web UI | HTTP |
| 9090 | Hub gRPC (Agent 通信) | gRPC |
| 9101 | Agent HTTP (metrics + health) | HTTP |
| 5432 | PostgreSQL | TCP |

## 部署清单

- [ ] PostgreSQL 数据库初始化
- [ ] Hub 二进制部署 + systemd 守护
- [ ] 每个 NGINX 节点部署 Agent + systemd 守护
- [ ] Prometheus 抓取 Hub + Agent /metrics
- [ ] 防火墙开放 3000（管理端口）和 9090（Agent 端口）

## 安全建议

1. Hub 的 3000 端口应仅对管理网络开放
2. JWT_SIGNING_KEY 使用固定随机字符串，不要用 ephemeral
3. Agent 的 9101 端口应对 Prometheus 开放
4. PostgreSQL 连接使用 `sslmode=require`（生产环境）