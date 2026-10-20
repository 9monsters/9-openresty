## ADR-004: PostgreSQL for Hub Storage

**Status:** Accepted  
**Date:** 2026-10-17

### Context

SQLite served the single-machine prototype well, but a multi-machine cluster hub has
different requirements: concurrent API requests during config fan-out, richer query needs
for node inventory and audit logs, and operator tooling (backup, replication).

### Decision

Replace SQLite with **PostgreSQL** on the hub. The agent remains stateless (no local DB).

### Schema considerations

The `collections` table uses a document-store pattern (JSON blobs keyed by collection name
and sequence). This maps naturally to PostgreSQL `JSONB`:

```sql
CREATE TABLE collections (
    collection TEXT NOT NULL,
    seq        INTEGER NOT NULL,
    id         BIGINT,
    data       JSONB NOT NULL,
    PRIMARY KEY (collection, seq)
);
CREATE INDEX idx_collections_id ON collections (collection, id);
```

`JSONB` gives us indexing on JSON paths if needed later, while keeping the flexible
schema the frontend already depends on.

### Connection

- `pgx/v5` with connection pool (`pgxpool`).
- Connection string from `--db-url` flag, env `DATABASE_URL`, or config file `database.url`.
- Default: `postgres://localhost:5432/openresty-admin?sslmode=disable`.

### Migration strategy

- `internal/store/migrations/` — numbered SQL files (`001_schema.sql`, …).
- `internal/store/migrate.go` — reads embedded migration files, tracks applied versions in a
  `schema_migrations` table.
- First migration creates the `collections` + `meta` tables matching the existing schema so
  existing SQLite data can be exported/imported.

### Consequences

- Deployment now requires a PostgreSQL instance (or docker compose for dev).
- The agent binary has no database dependency — it gets all state from the hub.
- `pgxpool` handles concurrent reads from multiple API goroutines without the single-writer
  bottleneck of the current `db.SetMaxOpenConns(1)`.
- Backup moves from file-copy to `pg_dump` or logical replication.