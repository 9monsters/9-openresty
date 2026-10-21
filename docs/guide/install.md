# 安装指南

## 环境要求

| 组件 | 最低版本 |
|------|---------|
| Go | >= 1.25 |
| PostgreSQL | >= 14 (或 SQLite) |
| Node.js | >= 22 (仅前端开发) |
| pnpm | >= 9 (仅前端开发) |
| OpenResty/NGINX | 被管理机器上需安装 |

## 快速开始

### 1. 数据库

**PostgreSQL**（默认）:
```bash
PGPASSWORD=your_password psql -h localhost -U postgres <<'SQL'
CREATE ROLE openresty LOGIN PASSWORD 'your_password';
CREATE DATABASE openresty OWNER openresty;
SQL
```

**SQLite**（轻量模式）:
```bash
export DATABASE_DRIVER=sqlite
export SQLITE_PATH=data/openresty-admin.db
```

### 2. 构建

```bash
pnpm --dir web build
go build -o bin/openresty-admin ./cmd/openresty-admin/
```

或一键构建：
```bash
pnpm build
```

### 3. 启动

```bash
JWT_SIGNING_KEY=my-secret-key ./bin/openresty-admin --port 3000
```

### 4. 访问

打开 [http://127.0.0.1:3000](http://127.0.0.1:3000)

默认管理员: `admin` / `admin123`

## 前端开发模式

```bash
# 终端1：启动 Go 后端
./bin/openresty-admin --port 3000 --frontend-mode=external

# 终端2：启动 Vue 开发服务器
pnpm --dir web dev
```

## 独立前端部署

```bash
STANDALONE=true pnpm --dir web build
# 产出在 web/dist/
```

## 下一步

- [配置文件](/reference/config-file) — JSON 配置文件
- [环境变量](/reference/env) — 支持的环境变量
- [CLI 参考](/reference/cli) — 命令行参数
