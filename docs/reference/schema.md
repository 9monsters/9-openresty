# 数据库表结构

## 认证

### users

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| username | TEXT UNIQUE NOT NULL | |
| email | TEXT NOT NULL | |
| password | TEXT NOT NULL | |
| role | TEXT DEFAULT 'Operator' | Administrator / Operator / Viewer |
| status | TEXT DEFAULT 'active' | |
| theme | TEXT DEFAULT 'auto' | |
| language | TEXT DEFAULT 'zh_CN' | |
| last_login | TIMESTAMPTZ | |
| created_at | TIMESTAMPTZ | |
| updated_at | TIMESTAMPTZ | |

### sessions

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| user_id | BIGINT FK | → users.id |
| token_jti | TEXT UNIQUE | JWT jti claim |
| expires_at | TIMESTAMPTZ | |
| created_at | TIMESTAMPTZ | |

## 节点

### nodes

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| name | TEXT NOT NULL | |
| hostname | TEXT NOT NULL | |
| ip | TEXT | |
| port | INT DEFAULT 9001 | |
| environment | TEXT DEFAULT 'production' | |
| nginx_version | TEXT | |
| is_openresty | BOOLEAN DEFAULT false | |
| status | TEXT DEFAULT 'offline' | online / offline / error |
| last_seen | TIMESTAMPTZ | |
| cert_fingerprint | TEXT | mTLS 证书指纹 |
| created_at / updated_at | TIMESTAMPTZ | |

## 站点

### sites

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| name | TEXT NOT NULL | |
| file_path | TEXT NOT NULL | |
| server_names | TEXT[] DEFAULT '{}' | |
| listens | TEXT[] DEFAULT '{80}' | |
| ssl | BOOLEAN DEFAULT false | |
| http2 | BOOLEAN DEFAULT false | |
| root_dir | TEXT | |
| access_log | TEXT | |
| error_log | TEXT | |
| proxy_pass | TEXT | |
| content | TEXT DEFAULT '' | 完整 server block |
| node_id | BIGINT FK | → nodes.id |
| enabled | BOOLEAN DEFAULT true | |
| created_at / updated_at | TIMESTAMPTZ | |

### site_locations

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| site_id | BIGINT FK | → sites.id |
| path | TEXT DEFAULT '/' | |
| proxy_pass | TEXT | |
| root_dir | TEXT | |
| return_directive | TEXT | |
| seq | INT DEFAULT 0 | |

### site_lua_directives

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| site_id | BIGINT FK | → sites.id |
| location_id | BIGINT FK | → site_locations.id |
| directive | TEXT NOT NULL | |
| value | TEXT NOT NULL | |

## 代理

### streams

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| name | TEXT NOT NULL | |
| listen_addr | TEXT NOT NULL | |
| proxy_pass | TEXT | |
| protocol | TEXT DEFAULT 'tcp' | |
| ssl | BOOLEAN DEFAULT false | |
| enabled | BOOLEAN DEFAULT true | |
| node_id | BIGINT FK | → nodes.id |
| created_at / updated_at | TIMESTAMPTZ | |

### upstreams

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| name | TEXT NOT NULL | |
| strategy | TEXT DEFAULT 'round_robin' | |
| node_id | BIGINT FK | → nodes.id |
| created_at / updated_at | TIMESTAMPTZ | |

### upstream_servers

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGSERIAL PK | |
| upstream_id | BIGINT FK | → upstreams.id |
| address | TEXT NOT NULL | |
| port | INT DEFAULT 80 | |
| weight | INT DEFAULT 1 | |
| max_fails | INT DEFAULT 3 | |
| fail_timeout | TEXT DEFAULT '30s' | |
| status | TEXT DEFAULT 'up' | |
| seq | INT DEFAULT 0 | |

## 证书

### certificates

| 列 | 类型 |
|----|------|
| id | BIGSERIAL PK |
| name, domains(TEXT[]), issuer, status, auto_renew, cert_pem, key_pem | |
| not_before, not_after | TIMESTAMPTZ |
| created_at, updated_at | TIMESTAMPTZ |

### acme_users

| 列 | 类型 |
|----|------|
| id | BIGSERIAL PK |
| name, email, provider, server_url, status | TEXT |
| cert_count | INT DEFAULT 0 |
| created_at, updated_at | TIMESTAMPTZ |

## 其他模块

### snippets / namespaces

`snippets`: id, name, content, description, is_template, created_at, updated_at
`namespaces`: id, name(UNIQUE), description, created_at, updated_at

### access_lists / access_list_rules

`access_lists`: id, name, description, created_at, updated_at
`access_list_rules`: id, list_id(FK), rule_type, ip, cidr, seq

### dns_credentials / dns_domains / dns_groups / dns_ddns

`dns_credentials`: id, name, provider, api_key, secret_key, created_at, updated_at
`dns_domains`: id, domain, credential_id(FK), record_count, created_at, updated_at
`dns_groups`: id, name, domains(TEXT[]), created_at, updated_at
`dns_ddns`: id, hostname, domain, ip, last_updated, created_at

### notifications / backups / auto_backup_settings

`notifications`: id, ntype, title, message, is_read, created_at
`backups`: id, name, backup_type, size, file_count, path, created_at
`auto_backup_settings`: id, enabled, interval, retention, schedule_time, path, updated_at

### config_history / config_releases / app_settings

`config_history`: id, file_path, content, message, sha256, created_at
`config_releases`: id, status, source, files(JSONB), node_results(JSONB), output, error, created_at, completed_at
`app_settings`: id, server_name, language, theme, nginx_bin, nginx_conf_dir, nginx_conf_path, nginx_log_dir, nginx_pid_path

### stats_system_history / stats_nginx_history / lua_scripts

`stats_system_history`: id, cpu_count, memory_used, memory_total, load_avg_1, load_avg_5, load_avg_15, uptime_secs, recorded_at
`stats_nginx_history`: id, active_connections, accepts, handled, requests, reading, writing, waiting, recorded_at
`lua_scripts`: id, name, directive, value, file_path, site_name, created_at, updated_at