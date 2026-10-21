# 二进制部署

## 构建

```bash
# 1. 构建前端（Vite 直接输出到 internal/embedweb/dist）
pnpm --dir packages/web build

# 2. 编译 Go 二进制
go build -o bin/openresty-admin ./cmd/openresty-admin/
```

产物：
- `bin/openresty-admin` — 单二进制，内嵌 Web UI

## 启动

```bash
JWT_SIGNING_KEY=prod-secret-key \
DATABASE_URL="postgres://openresty:pwd@localhost:5432/openresty?sslmode=disable" \
./bin/openresty-admin \
  --host 0.0.0.0 \
  --port 3000
```

SQLite 模式：

```bash
DATABASE_DRIVER=sqlite \
SQLITE_PATH=/var/lib/openresty-admin/data.db \
JWT_SIGNING_KEY=prod-secret-key \
./bin/openresty-admin --host 0.0.0.0 --port 3000
```

## 验证

```bash
# 健康检查
curl -s http://localhost:3000/api/system/info -H "Authorization: Bearer <token>"

# Prometheus 指标
curl -s http://localhost:3000/metrics
```

## 登录

```bash
TOKEN=$(curl -s -X POST http://localhost:3000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.token')
```
