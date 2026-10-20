# OpenResty Admin — Architecture Decision Records

## ADR-001: Multi-machine Cluster Deployment

**Status:** Accepted  
**Date:** 2026-10-17

### Context

The system manages OpenResty/NGINX instances across multiple machines from a single
management console.

### Decision

Adopt a **hub-agent** architecture:

- **Management Hub** — the Go binary running this repo. Serves the web UI, holds the SQLite
  store, coordinates config changes, and communicates with agents on remote nodes.
- **Node Agent** — a lightweight companion process on each managed host. Registers with the hub,
  reports health/stats, and applies config changes dispatched by the hub.

### Consequences

- Host Adapter (`internal/nginx/host.go`) must become an interface with two implementations:
  local (dev/single-machine) and remote (agent RPC).
- Store schema grows: the `nodes` collection already exists but was seed-only; it now becomes
  the authoritative node registry.
- Config release coordinator must fan out changes to target nodes and aggregate results.
- Terminal WebSocket must proxy to a remote shell.
- Network boundary introduces auth — API tokens or mutual TLS between hub and agents.

### Alternatives considered

- **SSH-based** (hub SSHs into nodes to run commands): simpler, no agent daemon, but slower
  fan-out and no real-time status push.
- **Ansible/puppet-style** (periodic pull): no real-time dashboard, no push-to-deploy.