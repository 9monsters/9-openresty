# CLI 命令参考

## openresty-admin

```
openresty-admin — OpenResty Admin Hub
Usage:
  openresty-admin [flags]
Flags:
  -c, --config string       JSON config file path
  -H, --host string          bind address (default 127.0.0.1)
  -p, --port int             bind port (default 3000)
  --grpc-addr string         gRPC listen address (default :9090)
  --nginx-bin string         path to openresty/nginx binary
  --nginx-conf-dir string    nginx config root
  --nginx-conf-path string   nginx main config file
  --nginx-log-dir string     nginx log directory
  --nginx-pid-path string    nginx PID file path
```

### 环境变量

| 变量 | 用途 |
|------|------|
| `DATABASE_URL` | PostgreSQL 连接串 |
| `HOST` | 监听地址 |
| `PORT` | HTTP 端口 |
| `JWT_SIGNING_KEY` | JWT 签名密钥（不设置则自动生成随机密钥） |

### 示例

```bash
# 最小化启动
JWT_SIGNING_KEY=mykey ./bin/openresty-admin

# 完整参数
JWT_SIGNING_KEY=prod-key \
DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=disable" \
./bin/openresty-admin \
  --host 0.0.0.0 \
  --port 3000 \
  --grpc-addr :9090 \
  --nginx-bin /usr/local/openresty/bin/openresty \
  --nginx-conf-dir /usr/local/etc/openresty

# 使用配置文件
./bin/openresty-admin -c /etc/openresty-admin/config.json
```

## openresty-agent

```
openresty-agent — OpenResty Node Agent
Usage:
  openresty-agent [flags]
Flags:
  --hub-addr string          hub gRPC address (default localhost:9090)
  --metrics-port int         metrics HTTP port (default 9101)
  --token string             registration token (hub-issued)
  --nginx-bin string         path to openresty/nginx binary
  --nginx-conf-dir string    nginx config root
  --nginx-conf-path string   nginx main config
  --nginx-log-dir string     nginx log directory
  --nginx-pid-path string    nginx PID file
```

### 示例

```bash
./bin/openresty-agent \
  --hub-addr hub.internal:9090 \
  --metrics-port 9101 \
  --nginx-bin /usr/local/openresty/bin/openresty
```