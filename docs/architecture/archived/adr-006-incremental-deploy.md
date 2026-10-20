## ADR-006: Incremental Config Deployment

**Status:** Accepted
**Date:** 2026-10-17

### Context

When the hub fans out a config change to multiple agents, transferring the entire
configuration tree each time is wasteful. A typical nginx config tree is hundreds of files
but only 1-3 change per release.

### Decision

Use **incremental diffs** — only the changed files are sent to each agent via gRPC.

### gRPC message

```protobuf
message PushConfigRequest {
  string release_id = 1;
  repeated ConfigChange changes = 2;  // only changed files
}

message ConfigChange {
  string path = 1;             // relative path within config root
  string content = 2;          // new content (empty for delete)
  bool delete = 3;             // remove the file
  string checksum_sha256 = 4;  // pre-validated hash from hub
}
```

This maps directly to `nginx.Change` in `internal/nginx/host.go` — the same struct the
release coordinator already uses. No translation layer needed.

### Flow

1. Hub computes changes for the release (same as single-machine mode).
2. Hub fans out `PushConfigRequest` to each target agent in parallel.
3. Each agent snapshots originals locally, writes the changed files.
4. Agent runs `nginx -t` on its own host.
5. On success: agent reloads nginx and responds `{ok: true}`.
6. On any agent failure: hub fans out rollback (`PushConfig` with original content) to all
   agents that succeeded.

### Local release journal

Each agent maintains a release journal at `<configDir>/.releases/<releaseID>/` containing:

```
.releases/
  42/
    manifest.json     -- {releaseID, source, files[], timestamp}
    files/
      servers/example.conf.orig  -- pre-change snapshot
      servers/example.conf.new   -- target content (for checksum verification)
```

This is separate from `.backups/` which is user-initiated manual snapshots shown in the UI.

### Rollback

1. Hub detects agent failure (gRPC error or explicit `{ok: false}`).
2. Hub sends `PushConfig` with the original file contents to each agent in the batch.
3. Agent restores from its local journal `.releases/<id>/files/*.orig`.
4. Agent runs `nginx -t` after restore — never reload an untested config, even on rollback.
5. If rollback succeeds: agent responds `{ok: true, rolled_back: true}`.
6. If rollback fails: agent enters `recovery_required` state — human intervention needed.

### Dedup across nodes

If multiple nodes share identical config content (e.g., a template `proxy.conf`), the hub
sends the same change payload to each. No content-addressed dedup at the gRPC layer — HTTP/2
header compression is sufficient for the scale (dozens of nodes, not thousands).

### Consequences

- `nginx.Change` struct is already designed for this — no schema change in the hub.
- Agent needs a local release journal directory for snapshot/rollback.
- Fan-out latency is bounded by the slowest agent — all agents must validate before any
  reload. The coordinator extends naturally from serial-local to parallel-remote.
- No Git dependency on the agent — the push model keeps the agent simple.