# 二进制部署

## 构建

```bash
# 1. 构建前端
pnpm --filter @openresty-admin/web build

# 2. 拷贝到嵌入目录
bash scripts/copy-web-dist.sh

# 3. 编译 Go 二进制
go build -o bin/openresty-admin ./cmd/openresty-admin/
go build -o bin/openresty-agent ./cmd/openresty-agent/
```

产物：
- `bin/openresty-admin` (24MB) — Hub + Web UI
- `bin/openresty-agent` (14MB) — Agent 守护进程

## Hub 启动

```bash
JWT_SIGNING_KEY=prod-secret-key \
DATABASE_URL="postgres://openresty:pwd@db.internal:5432/openresty?sslmode=disable" \
./bin/openresty-admin \
  --host 0.0.0.0 \
  --port 3000 \
  --grpc-addr :9090
```

## Agent 启动

在每个受控节点上：

```bash
./bin/openresty-agent \
  --hub-addr hub.internal:9090 \
  --metrics-port 9101 \
  --nginx-bin /usr/local/bin/openresty \
  --nginx-conf-dir /etc/nginx \
  --nginx-log-dir /var/log/nginx
```

## 验证

```bash
# Hub 健康检查
curl -s http://hub:3000/api/system/info -H "Authorization: Bearer <token>"

# Agent 健康检查
curl -s http://node:9101/health

# Prometheus 指标
curl -s http://hub:3000/metrics
curl -s http://node:9101/metrics
```

## 登录

```bash
TOKEN=$(curl -s -X POST http://hub:3000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')
```