## ADR-005: Fat Agent Architecture

**Status:** Accepted
**Date:** 2026-10-17

### Context

Each managed node runs an `openresty-agent` binary. The question is whether the agent is a
thin proxy (hub does all the work) or a capable local manager (agent owns its nginx ops).

### Decision

Use a **fat agent** model. The agent exposes the same Host Adapter interface that
`internal/nginx/host.go` provides today — config file I/O, nginx binary operations,
backup/restore, site parsing. The hub becomes an **orchestration layer**:

- Agent registers with the hub, reports its identity and capabilities.
- Hub fans out commands (validate, write, reload) but each agent executes them locally.
- Agent runs its own `nginx -t` before reload — the hub never runs `nginx -t` remotely.
- Agent maintains its own config backups (`.backups/` directory locally).
- Each agent can be operated directly on the node without the hub (useful for debugging
  and emergency recovery).

### Hub responsibilities (what changes)

| Before (single-machine) | After (multi-machine with fat agents) |
|--------------------------|---------------------------------------|
| `Host.WriteConfig()` directly | gRPC call to agent's `ConfigDeploy.Write` |
| `Host.Test()` directly | gRPC call to agent's `ConfigDeploy.Validate` |
| `Host.Reload()` directly | gRPC call to agent's `ConfigDeploy.Reload` |
| `Host.Stats()` curl localhost | Agent streams stats to hub via `StatsCollector.StreamStats` |
| Terminal spawns local shell | Hub proxies WebSocket to agent's `TerminalProxy.Shell` |
| Single `config-releases` collection | Hub's release becomes a fan-out: one hub record + per-agent results |

### Agent internal structure

```
cmd/openresty-agent/
  main.go          -- cobra CLI: --hub-addr, --token, --cert, --key
  server.go        -- gRPC server + health loop

internal/
  agent/
    services.go    -- ConfigDeployService, StatsService, TerminalService, LogService
  nginx/           -- reuses the existing Host adapter (host.go, parser.go, backup.go)
  config/          -- reuses config.go (paths, defaults)
```

### State ownership

- **Hub owns**: users, nodes registry, notifications, DNS credentials, ACME users,
  snippets, access lists, upstreams, streams, certificates — all `collections` that
  are management-plane data.
- **Agent owns**: nginx config files on its host, local backups, local log files.
- **Hub caches**: agent health, nginx status, stats — refreshed via streaming/polling.

### Release coordinator changes

The current `internal/release/coordinator.go` serializes releases for a single host.
The multi-node coordinator:

1. Hub creates a `config-releases` record with status `running`.
2. Hub fans out `Validate` to each target agent (parallel).
3. On all-agents-valid, hub fans out `Write` + `Reload` to each agent.
4. Any agent failure triggers rollback on the agents that succeeded.
5. Hub records per-agent results, sets the aggregate release status.

### Consequences

- Agent binary is a full Go binary (~15-20MB) with the nginx parser, backup, and store
  packages compiled in. Acceptable for server deployments.
- No central bottleneck on config validation — each node validates its own config locally.
- Terminal sessions connect directly to the target node's shell without hub-side nginx
  knowledge.
- Agent can be developed and tested standalone: `openresty-agent serve --dev` runs a local
  instance with `grpcurl`-testable endpoints.