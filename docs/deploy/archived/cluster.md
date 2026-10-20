# 集群部署

## 架构

```
         ┌─────────────────┐
         │   Hub (:3000)   │  Web UI + gRPC Server + PostgreSQL
         └────────┬────────┘
                  │ gRPC (:9090)
     ┌────────────┼────────────┐
     ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Agent 1 │ │ Agent 2 │ │ Agent 3 │  各节点运行 openresty-agent
│ nginx   │ │ nginx   │ │ nginx   │
└─────────┘ └─────────┘ └─────────┘
```

## Hub 部署

1. 部署 PostgreSQL
2. 初始化数据库
3. 启动 Hub 服务

## Agent 部署

在每个受控节点：

1. 安装 OpenResty/NGINX
2. 部署 `openresty-agent` 二进制
3. 配置 systemd 守护
4. Agent 自动注册到 Hub

## 配置下发流程

```
用户操作 (Web UI)
  → Hub 创建 Release 记录
  → Hub 并行 gRPC PushConfig 到目标 Agent
  → Agent 写入文件 → nginx -t 本地验证 → reload
  → 任一节点失败 → Hub 扇出 Rollback
  → 全部成功 → 记录 Release 为 succeeded
```

## Prometheus 监控

```yaml
scrape_configs:
  - job_name: openresty-hub
    static_configs: [{ targets: ['hub:3000'] }]
  - job_name: openresty-agents
    static_configs: [{ targets: ['node-01:9101', 'node-02:9101', 'node-03:9101'] }]
```