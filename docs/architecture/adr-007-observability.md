## ADR-007: Structured Logging, Distributed Tracing, and Prometheus Metrics

**Status:** Accepted
**Date:** 2026-10-17

### Context

In a multi-machine deployment, a single user action (e.g., "deploy site config") fans out
across hub → multiple agents. The operator needs three layers of observability:

1. **Logs** — per-request detail, searchable by trace ID across hub and agents.
2. **Traces** — end-to-end request flow from browser to agent and back.
3. **Metrics** — aggregate health of the cluster (agent online count, deploy success rate,
   API latency percentiles).

### Decision

**A — Structured logging (`log/slog`)** with JSON output on both hub and agent.
**B — Distributed trace IDs** propagated through HTTP headers and gRPC metadata.
**C — Prometheus metrics** at `GET /metrics` on both hub and agent.

### Slog configuration

```go
// hub and agent both use this
handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
```

Log lines are one-line JSON objects: `time`, `level`, `msg`, plus contextual attributes
including `trace_id`.

### Trace propagation

- Hub generates a `trace_id` (UUIDv7) per incoming HTTP request via middleware.
- gRPC unary/stream interceptors on both sides propagate `trace_id` via gRPC metadata
  (`x-trace-id`).
- Agent attaches `trace_id` to every log line it emits while handling that request.
- Browser can optionally send `X-Trace-Id` header; if absent, hub generates one.

```
# Hub log (inbound HTTP request)
{"time":"...","level":"INFO","msg":"request","trace_id":"018...","method":"POST","path":"/api/sites","user":"admin"}

# Hub log (outbound gRPC fan-out)
{"time":"...","level":"INFO","msg":"config:push","trace_id":"018...","agent":"node-3","release_id":42,"files":2}

# Agent log (inbound gRPC)
{"time":"...","level":"INFO","msg":"config:write","trace_id":"018...","path":"sites-enabled/example.conf"}

# Agent log (nginx reload)
{"time":"...","level":"INFO","msg":"nginx:reload","trace_id":"018...","ok":true,"duration_ms":342}

# Hub log (response)
{"time":"...","level":"INFO","msg":"response","trace_id":"018...","status":200,"duration_ms":1423}
```

### Prometheus metrics

Registered on both hub and agent at `GET /metrics` using `prometheus/client_golang`.

#### Hub metrics

| Metric | Type | Description |
|--------|------|-------------|
| `openresty_admin_http_requests_total` | Counter | API requests by method, path, status |
| `openresty_admin_http_request_duration_seconds` | Histogram | API request latency (p50/p95/p99) |
| `openresty_admin_agents_online` | Gauge | Number of connected agents |
| `openresty_admin_deploy_total` | Counter | Config releases by status (success/failed/rolled_back) |
| `openresty_admin_deploy_duration_seconds` | Histogram | End-to-end deploy latency per release |
| `openresty_admin_agent_rpc_duration_seconds` | Histogram | gRPC call latency per agent, per method |
| `openresty_admin_agent_rpc_errors_total` | Counter | gRPC errors per agent, per method |

#### Agent metrics

| Metric | Type | Description |
|--------|------|-------------|
| `openresty_agent_nginx_status` | Gauge | 1=running, 0=stopped, -1=unknown |
| `openresty_agent_nginx_config_test_total` | Counter | `nginx -t` calls by result (ok/fail) |
| `openresty_agent_nginx_reload_total` | Counter | nginx reloads by result (ok/fail) |
| `openresty_agent_deploy_total` | Counter | Received deploys by status |
| `openresty_agent_deploy_duration_seconds` | Histogram | Per-agent deploy latency |
| `openresty_agent_connections_active` | Gauge | Active nginx connections (from stub_status) |
| `openresty_agent_connections_total` | Counter | Total nginx connections handled |

### Metrics endpoint

```go
// internal/metrics/hub.go
func RegisterHubMetrics() http.Handler {
    registry := prometheus.NewRegistry()
    // ... register all hub metrics ...
    return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}

// internal/metrics/agent.go
func RegisterAgentMetrics() http.Handler {
    registry := prometheus.NewRegistry()
    // ... register all agent metrics ...
    return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}
```

Both hub and agent serve `/metrics` on their HTTP listener (the hub on its main port, the
agent on a separate management port, default `9101`).

### gRPC interceptors (trace + metrics combined)

The trace interceptor and the metrics interceptor are separate middleware layers:

```go
// internal/trace/interceptor.go
func UnaryClientInterceptor() grpc.UnaryClientInterceptor { /* trace_id propagation */ }
func UnaryServerInterceptor() grpc.UnaryServerInterceptor { /* trace_id extraction */ }

// internal/metrics/grpc_interceptor.go
func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
    // records openresty_admin_agent_rpc_duration_seconds and _errors_total
}
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
    // records openresty_agent_deploy_total and deploy_duration_seconds
}
```

### Logs are local, metrics are scraped

- Agent operational logs stay on the agent host — the hub does not collect them centrally.
  Operators SSH into the agent or point a log shipper at the agent's stdout.
- Prometheus scrapes `/metrics` from both hub and agent. A `prometheus.yml` example:

```yaml
scrape_configs:
  - job_name: openresty-hub
    static_configs:
      - targets: ['hub:3000']
  - job_name: openresty-agents
    static_configs:
      - targets: ['node-01:9101', 'node-02:9101', 'node-03:9101']
```

### Consequences

- No external tracing dependency (no Jaeger/Tempo required). Traces are reconstructed by
  grepping log files for a `trace_id`. OpenTelemetry can be added later as an export layer
  on the existing `trace_id` plumbing.
- Prometheus metrics give operators alertable signals (agents offline, deploy failures)
  without log diving.
- Agent requires an additional HTTP port (9101) for `/metrics` and health checks.
- `prometheus/client_golang` adds ~2MB to each binary — acceptable for server deployments.
- The existing frontend dashboard gains a metrics panel powered by Prometheus queries
  (future work — Phase 5+).