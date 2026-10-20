## ADR-008: Implementation Roadmap — Bottom-Up

**Status:** Accepted
**Date:** 2026-10-17

### Decision

Build the cluster architecture **bottom-up**: infrastructure layers first, then agent
communication, then fat agent wiring, then cluster deployment, and finally authentication.

### Phases

#### Phase 1 — Infrastructure foundation

- Replace SQLite with PostgreSQL via `pgx/v5` + `pgxpool` (ADR-004).
- `internal/store/migrations/` with embedded SQL files, `schema_migrations` tracker.
- Replace `fmt.Print` with `log/slog` JSON handler throughout hub (ADR-007).
- Add trace ID middleware: `internal/trace/http.go` (HTTP middleware), `internal/trace/context.go` (context key).
- Add Prometheus metrics: `internal/metrics/hub.go` — `/metrics` endpoint with HTTP request
  counters/histograms, agent online gauge, deploy counters (ADR-007).
- Add gRPC metrics interceptor: `internal/metrics/grpc_interceptor.go` — client-side RPC
  latency/error counters.
- `go.mod`: add `github.com/jackc/pgx/v5`, `github.com/prometheus/client_golang`;
  remove `modernc.org/sqlite`.

**Delivery**: hub compiles with PG, emits JSON logs with trace IDs, serves `/metrics`.
Tests pass against a real PG (test helper: `PGTEST_DSN` env var).

#### Phase 2 — gRPC contract and agent scaffold

- Define protobuf services in `proto/openresty/admin/v1/` (ADR-002):
  `AgentRegistration`, `ConfigDeploy`, `StatsCollector`, `TerminalProxy`, `LogReader`.
- Generate Go code via `buf`: `gen/go/openresty/admin/v1/`.
- Build `cmd/openresty-agent/` skeleton: cobra CLI with `--hub-addr`, `--cert`, `--key`.
- Agent serves its own `/metrics` on port 9101 (ADR-007).
- Implement gRPC server with stub handlers, gRPC client in hub.
- mTLS: hub generates CA, signs agent certs. `nodes` collection stores cert fingerprints.

**Delivery**: `openresty-agent serve` starts, connects to hub, registers. Hub lists
connected agents. Prometheus scrapes both hub and agent /metrics.

#### Phase 3 — Fat agent (Host Adapter reuse)

- Define `Host` interface in `internal/nginx/host.go` from existing method set.
- `LocalHost` — the current implementation (single-machine, used by agent).
- `RemoteHost` — gRPC client that proxies to the agent (used by hub).
- Agent wires its gRPC services to `LocalHost`.
- Terminal proxy: hub WebSocket → `TerminalProxy.Shell` bidirectional stream → agent shell.
- Stats streaming: `StatsCollector.StreamStats` pushes nginx connection stats to hub.

**Delivery**: hub API routes work via agent gRPC. Terminal connects to remote agent shell.
Stats stream from agent to hub, exposed in dashboard. All traces carry `trace_id` across
hub↔agent boundary.

#### Phase 4 — Incremental config deployment (ADR-006)

- Hub side: release coordinator fans out `PushConfigRequest` to target agents in parallel.
- Agent side: `ConfigDeploy` service applies incremental changes, snapshots originals to
  local `.releases/<id>/` journal, runs `nginx -t`, reloads, reports result.
- Release coordinator rewrite: `ParallelPublish(changes, nodeIDs)` with per-agent status
  aggregation. All-agents-valid before any reload. Failure triggers rollback to
  succeeded agents.
- `config-releases` collection gains a `node_results` JSONB array: per-agent status,
  output, and duration.

**Delivery**: operator creates a site in the UI → hub fans out to agents in parallel →
each agent validates and reloads locally → aggregate result recorded with per-node detail.

#### Phase 5 — Authentication

- JWT auth middleware (ADR-003): `internal/auth/jwt.go`, `internal/auth/middleware.go`.
- `POST /api/auth/login`, `POST /api/auth/refresh`.
- bcrypt migration: transparent upgrade of plaintext passwords on first login.
- Hub-agent mTLS enforcement: hub refuses unauthenticated agent connections.
- Frontend: login page, token storage in Pinia, axios interceptor attaches `Authorization` header.

**Delivery**: API returns 401 without token. Login flow works end-to-end. Agents
authenticated via mTLS.

### Dependency graph

```
Phase 1 (PG + slog + trace + Prometheus)
  └── Phase 2 (protobuf + agent scaffold + agent /metrics)
        └── Phase 3 (Host interface, fat agent, terminal, stats streaming)
              └── Phase 4 (incremental push deploy, parallel coordinator)
                    └── Phase 5 (JWT, mTLS enforcement)
```

Each phase is independently buildable and testable. Phase 2 depends on Phase 1's database
and observability infra. Phase 5 can begin in parallel with Phase 4 once Phase 3 lands.

### What stays unchanged

- Frontend (Vue 3, antdv-next, UnoCSS) — no architectural changes, only new routes/pages
  added as needed (login page, node list with health indicators, metrics dashboard).
- All existing API route shapes — the JSON contracts between frontend and hub remain
  compatible. gRPC is hub↔agent only, not browser↔hub.
- `internal/nginx/parser.go` and `internal/nginx/backup.go` — used as-is by the agent.
- `internal/api/*.go` — handlers stay the same; only the `Host` interface they call changes
  from local to remote.