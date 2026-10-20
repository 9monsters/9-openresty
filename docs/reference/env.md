# 环境变量

以下环境变量用于启动时配置。

## 必需

| 变量 | 说明 | 示例 |
|------|------|------|
| `JWT_SIGNING_KEY` | JWT 签名密钥 | `dev-secret-key` |
| `PG_DSN` | PostgreSQL 连接字符串 | `host=192.168.5.100 user=openresty password=... dbname=openresty sslmode=disable` |

## 可选

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `PORT` | HTTP 监听端口 | `3000` |
| `HOST` | 绑定地址 | `127.0.0.1` |
| `GRPC_PORT` | gRPC 端口（Hub 模式） | `9090` |
| `CONFIG_FILE` | 配置文件路径 | 无 |
| `LOG_LEVEL` | 日志级别 | `info` |
| `LOG_FORMAT` | 日志格式（json/text） | `text` |
| `METRICS_ENABLED` | 启用 Prometheus /metrics | `true` |
| `METRICS_PORT` | Metrics 端口 | `9091` |

## 优先级

```
CLI 参数 > 配置文件 > 环境变量 > 默认值
```