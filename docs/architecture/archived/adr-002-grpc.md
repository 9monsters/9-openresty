## ADR-002: Hub-Agent Communication via gRPC

**Status:** Accepted  
**Date:** 2026-10-17

### Context

The hub must push config changes, proxy terminal sessions, and collect health/stats from
multiple remote agents. The protocol must be typed, bidirectional, and handle streaming.

### Decision

Use **gRPC** with protobuf service definitions.

Reasons:
- **Bidirectional streaming** is first-class — the terminal proxy and health push need it.
- **Code generation** eliminates hand-rolled serialization bugs across Go services.
- **HTTP/2 multiplexing** means one persistent connection per agent, no connection pool noise.
- The ecosystem is Go-native; both hub and agent are written in Go.

### Service boundaries (candidate)

| Service | RPCs | Direction |
|---------|------|-----------|
| `AgentRegistration` | `Register`, `Heartbeat` | agent → hub |
| `ConfigDeploy` | `PushConfig`, `ValidateConfig`, `Reload` | hub → agent |
| `StatsCollector` | `StreamStats` | agent → hub (server-streaming) |
| `TerminalProxy` | `Shell` | bidirectional stream |
| `LogReader` | `ReadLog` | hub → agent (server-streaming) |

### Consequences

- Build pipeline requires `protoc` + `protoc-gen-go` + `protoc-gen-go-grpc`.
- `internal/nginx/host.go` splits into `Host` interface + `local` / `remote` implementations.
- The `remote` host adapter calls the agent via gRPC instead of local `os/exec`.
- Agent must ship as a separate binary (`cmd/openresty-agent/`) with its own cobra CLI.
- Mutual TLS between hub and agents — agent certs signed by the hub's CA, registered
  fingerprints stored in the `nodes` collection.