# 快速开始

OpenResty Admin 是一个 Web-based 的 OpenResty/NGINX 集群管理平台。

## 系统要求

| 组件 | 要求 |
|------|------|
| Go | >= 1.25 |
| Node.js | >= 22 (仅前端开发) |
| PostgreSQL | >= 14 |
| pnpm | >= 9 (仅前端开发) |
| OpenResty/NGINX | 被管理节点上需安装 |

## 5 分钟体验

```bash
# 1. 初始化数据库
PGPASSWORD=1aac@Postgres psql -h 192.168.5.100 -U postgres <<'SQL'
CREATE ROLE openresty LOGIN PASSWORD '1aac@Openresty';
CREATE DATABASE openresty OWNER openresty;
SQL

# 2. 构建（含前端）
pnpm --filter @openresty-admin/web build
bash scripts/copy-web-dist.sh
go build -o bin/openresty-admin ./cmd/openresty-admin/

# 3. 启动
JWT_SIGNING_KEY=my-secret-key ./bin/openresty-admin --port 3000

# 4. 访问
open http://127.0.0.1:3000
```

默认管理员账号: `admin` / `admin123`

## 下一步

- [安装指南](/guide/install) - 详细安装步骤
- [配置文件](/guide/config) - JSON 配置文件说明
- [部署指南](/deploy/) - 生产环境部署