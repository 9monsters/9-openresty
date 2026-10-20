## ADR-009: Normalized Database Schema

**Status:** Accepted
**Date:** 2026-10-18

### Decision

Replace the generic JSONB `collections` table with normalized, typed tables. Every
module owns its own table(s) with proper column types, constraints, and indexes.

### Schema (25 tables)

#### 1. auth — users & sessions

```sql
CREATE TABLE users (
  id           BIGSERIAL PRIMARY KEY,
  username     TEXT NOT NULL UNIQUE,
  email        TEXT NOT NULL,
  password     TEXT NOT NULL,           -- bcrypt hash
  role         TEXT NOT NULL DEFAULT 'Operator',  -- Administrator | Operator | Viewer
  status       TEXT NOT NULL DEFAULT 'active',
  theme        TEXT NOT NULL DEFAULT 'auto',
  language     TEXT NOT NULL DEFAULT 'zh_CN',
  last_login   TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
  id           BIGSERIAL PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_jti    TEXT NOT NULL UNIQUE,    -- JWT jti claim for revocation
  expires_at   TIMESTAMPTZ NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 2. nodes — agent-managed hosts

```sql
CREATE TABLE nodes (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  hostname      TEXT NOT NULL,
  ip            TEXT,
  port          INT NOT NULL DEFAULT 9001,
  environment   TEXT NOT NULL DEFAULT 'production',
  nginx_version TEXT,
  is_openresty  BOOLEAN NOT NULL DEFAULT false,
  status        TEXT NOT NULL DEFAULT 'offline',  -- online | offline | error
  last_seen     TIMESTAMPTZ,
  labels        JSONB DEFAULT '{}',
  cert_fingerprint TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 3. sites — server blocks

```sql
CREATE TABLE sites (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  file_path     TEXT NOT NULL,
  server_names  TEXT[] NOT NULL DEFAULT '{}',
  listens       TEXT[] NOT NULL DEFAULT '{80}',
  ssl           BOOLEAN NOT NULL DEFAULT false,
  http2         BOOLEAN NOT NULL DEFAULT false,
  root          TEXT,
  access_log    TEXT,
  error_log     TEXT,
  proxy_pass    TEXT,
  content       TEXT NOT NULL,           -- full server block text
  node_id       BIGINT REFERENCES nodes(id),
  enabled       BOOLEAN NOT NULL DEFAULT true,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE site_locations (
  id            BIGSERIAL PRIMARY KEY,
  site_id       BIGINT NOT NULL REFERENCES sites(id) ON DELETE CASCADE,
  path          TEXT NOT NULL DEFAULT '/',
  proxy_pass    TEXT,
  root          TEXT,
  return_directive TEXT,
  seq           INT NOT NULL DEFAULT 0
);

CREATE TABLE site_lua_directives (
  id            BIGSERIAL PRIMARY KEY,
  site_id       BIGINT REFERENCES sites(id) ON DELETE CASCADE,
  location_id   BIGINT REFERENCES site_locations(id) ON DELETE CASCADE,
  directive     TEXT NOT NULL,           -- e.g. content_by_lua_block
  value         TEXT NOT NULL,
  node_id       BIGINT REFERENCES nodes(id)
);
```

#### 4. streams

```sql
CREATE TABLE streams (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  listen        TEXT NOT NULL,
  proxy_pass    TEXT,
  protocol      TEXT NOT NULL DEFAULT 'tcp',
  ssl           BOOLEAN NOT NULL DEFAULT false,
  enabled       BOOLEAN NOT NULL DEFAULT true,
  node_id       BIGINT REFERENCES nodes(id),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 5. upstreams

```sql
CREATE TABLE upstreams (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  strategy      TEXT NOT NULL DEFAULT 'round_robin',
  node_id       BIGINT REFERENCES nodes(id),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE upstream_servers (
  id            BIGSERIAL PRIMARY KEY,
  upstream_id   BIGINT NOT NULL REFERENCES upstreams(id) ON DELETE CASCADE,
  address       TEXT NOT NULL,
  port          INT NOT NULL DEFAULT 80,
  weight        INT NOT NULL DEFAULT 1,
  max_fails     INT NOT NULL DEFAULT 3,
  fail_timeout  TEXT NOT NULL DEFAULT '30s',
  status        TEXT NOT NULL DEFAULT 'up',
  seq           INT NOT NULL DEFAULT 0
);
```

#### 6. certificates & acme

```sql
CREATE TABLE certificates (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  domains       TEXT[] NOT NULL DEFAULT '{}',
  issuer        TEXT NOT NULL DEFAULT 'Manual',
  status        TEXT NOT NULL DEFAULT 'valid',
  auto_renew    BOOLEAN NOT NULL DEFAULT false,
  cert_pem      TEXT NOT NULL DEFAULT '',
  key_pem       TEXT NOT NULL DEFAULT '',
  not_before    TIMESTAMPTZ,
  not_after     TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE acme_users (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  email         TEXT NOT NULL,
  provider      TEXT NOT NULL DEFAULT 'Let''s Encrypt',
  server_url    TEXT NOT NULL DEFAULT 'https://acme-v02.api.letsencrypt.org/directory',
  status        TEXT NOT NULL DEFAULT 'active',
  cert_count    INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 7. access_lists

```sql
CREATE TABLE access_lists (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  description   TEXT NOT NULL DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE access_list_rules (
  id            BIGSERIAL PRIMARY KEY,
  list_id       BIGINT NOT NULL REFERENCES access_lists(id) ON DELETE CASCADE,
  rule_type     TEXT NOT NULL DEFAULT 'allow',  -- allow | deny
  ip            TEXT NOT NULL,
  cidr          INT,
  seq           INT NOT NULL DEFAULT 0
);
```

#### 8. snippets

```sql
CREATE TABLE snippets (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  content       TEXT NOT NULL DEFAULT '',
  description   TEXT NOT NULL DEFAULT '',
  is_template   BOOLEAN NOT NULL DEFAULT false,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 9. dns — 4 tables

```sql
CREATE TABLE dns_credentials (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  provider      TEXT NOT NULL,
  api_key       TEXT NOT NULL,
  secret_key    TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dns_domains (
  id            BIGSERIAL PRIMARY KEY,
  domain        TEXT NOT NULL,
  credential_id BIGINT REFERENCES dns_credentials(id),
  record_count  INT NOT NULL DEFAULT 0,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dns_groups (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  domains       TEXT[] NOT NULL DEFAULT '{}',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE dns_ddns (
  id            BIGSERIAL PRIMARY KEY,
  hostname      TEXT NOT NULL,
  domain        TEXT NOT NULL,
  ip            TEXT,
  last_updated  TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 10. notifications

```sql
CREATE TABLE notifications (
  id            BIGSERIAL PRIMARY KEY,
  type          TEXT NOT NULL DEFAULT 'info',  -- success | warning | error | info
  title         TEXT NOT NULL,
  message       TEXT NOT NULL,
  is_read       BOOLEAN NOT NULL DEFAULT false,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 11. backups

```sql
CREATE TABLE backups (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  backup_type   TEXT NOT NULL DEFAULT 'manual',  -- auto | manual
  size          BIGINT NOT NULL DEFAULT 0,
  file_count    INT NOT NULL DEFAULT 0,
  path          TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE auto_backup_settings (
  id            BIGSERIAL PRIMARY KEY DEFAULT 1,
  enabled       BOOLEAN NOT NULL DEFAULT true,
  interval      TEXT NOT NULL DEFAULT 'daily',   -- hourly | daily | weekly
  retention     INT NOT NULL DEFAULT 7,
  schedule_time TEXT NOT NULL DEFAULT '03:00',
  path          TEXT NOT NULL DEFAULT '/var/backups/openresty',
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Previous config snapshots, replacing filesystem .backups/
CREATE TABLE config_history (
  id            BIGSERIAL PRIMARY KEY,
  file_path     TEXT NOT NULL,
  content       TEXT NOT NULL,
  message       TEXT NOT NULL DEFAULT '',
  sha256        TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_config_history_path ON config_history(file_path);
```

#### 12. config_releases

```sql
CREATE TABLE config_releases (
  id            BIGSERIAL PRIMARY KEY,
  release_id    BIGINT,
  status        TEXT NOT NULL DEFAULT 'running',
  source        TEXT NOT NULL DEFAULT 'config',
  files         JSONB NOT NULL DEFAULT '[]',
  node_results  JSONB,
  output        TEXT,
  error         TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at  TIMESTAMPTZ
);
```

#### 13. app_settings

```sql
CREATE TABLE app_settings (
  id              BIGSERIAL PRIMARY KEY DEFAULT 1,
  server_name     TEXT NOT NULL DEFAULT 'OpenResty Admin',
  language        TEXT NOT NULL DEFAULT 'zh_CN',
  theme           TEXT NOT NULL DEFAULT 'auto',
  nginx_bin       TEXT NOT NULL,
  nginx_conf_dir  TEXT NOT NULL,
  nginx_conf_path TEXT NOT NULL,
  nginx_log_dir   TEXT NOT NULL,
  nginx_pid_path  TEXT NOT NULL,
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 14. namespaces

```sql
CREATE TABLE namespaces (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL UNIQUE,
  description   TEXT DEFAULT '',
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### 15. stats history (time-series)

```sql
CREATE TABLE stats_system_history (
  id            BIGSERIAL PRIMARY KEY,
  cpu_count     INT NOT NULL,
  memory_used   BIGINT NOT NULL,
  memory_total  BIGINT NOT NULL,
  load_avg_1    DOUBLE PRECISION NOT NULL,
  load_avg_5    DOUBLE PRECISION NOT NULL,
  load_avg_15   DOUBLE PRECISION NOT NULL,
  uptime_secs   BIGINT NOT NULL,
  recorded_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_stats_system_time ON stats_system_history(recorded_at);

CREATE TABLE stats_nginx_history (
  id                  BIGSERIAL PRIMARY KEY,
  active_connections  BIGINT NOT NULL DEFAULT 0,
  accepts             BIGINT NOT NULL DEFAULT 0,
  handled             BIGINT NOT NULL DEFAULT 0,
  requests            BIGINT NOT NULL DEFAULT 0,
  reading             BIGINT NOT NULL DEFAULT 0,
  writing             BIGINT NOT NULL DEFAULT 0,
  waiting             BIGINT NOT NULL DEFAULT 0,
  recorded_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_stats_nginx_time ON stats_nginx_history(recorded_at);
```

#### 16. lua scripts

```sql
CREATE TABLE lua_scripts (
  id            BIGSERIAL PRIMARY KEY,
  name          TEXT NOT NULL,
  directive     TEXT NOT NULL,
  value         TEXT NOT NULL,
  file_path     TEXT,
  site_name     TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Migration strategy

1. Keep the `collections` table temporarily during migration.
2. New migration version creates all new tables.
3. Data migration function reads from `collections` → inserts into typed tables.
4. After verification, drop the `collections` table.

### Consequences

- Frontend API shapes remain compatible — each handler maps typed rows to the
  same JSON response shape the Vue views expect.
- `crudSpec` in `collections.go` becomes per-table handler registration with
  typed queries instead of generic JSONB load/save.
- Performance: indexed queries on `server_names`, `status`, `recorded_at`,
  `username` replace full-table JSONB scans.
- Stats retention: `stats_system_history` and `stats_nginx_history` keep
  per-second data for 7 days, then roll up to 1-minute buckets (future task).