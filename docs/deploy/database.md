# 数据库初始化

## 创建用户和数据库

```sql
-- 以 postgres 超级用户执行
CREATE ROLE openresty LOGIN PASSWORD '1aac@Openresty';
CREATE DATABASE openresty OWNER openresty;
GRANT ALL PRIVILEGES ON DATABASE openresty TO openresty;
```

或通过项目脚本：

```bash
PGPASSWORD=1aac@Postgres psql -h 192.168.5.100 -U postgres -f scripts/init-db.sql
```

## 自动迁移

服务启动时 GORM `AutoMigrate` 自动创建/更新 30+ 张表。无需手动执行 DDL。

首次启动会自动创建默认管理员 `admin / admin123`。

## 连接配置

通过环境变量指定：

```bash
export DATABASE_URL="postgres://openresty:1aac%40Openresty@192.168.5.100:5432/openresty?sslmode=disable"
```

密码中的特殊字符需 URL 编码（`@` → `%40`）。

## 表清单

| 分类 | 表名 |
|------|------|
| 认证 | `users`, `sessions` |
| 节点 | `nodes` |
| 站点 | `sites`, `site_locations`, `site_lua_directives` |
| 代理 | `streams`, `upstreams`, `upstream_servers` |
| 证书 | `certificates`, `acme_users` |
| ACL | `access_lists`, `access_list_rules` |
| DNS | `dns_credentials`, `dns_domains`, `dns_groups`, `dns_ddns` |
| 其他 | `snippets`, `namespaces`, `notifications` |
| 备份 | `backups`, `auto_backup_settings`, `config_history` |
| 设置 | `app_settings`, `config_releases` |
| 统计 | `stats_system_history`, `stats_nginx_history` |
| Lua | `lua_scripts` |