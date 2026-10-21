# OpenResty Admin — Agent Guidelines

Web-based OpenResty/NGINX management interface: Go backend (cobra CLI + embedded Vue 3 frontend).

## Environment boundary

- **Backend**: Go with cobra CLI, serves on port 3000 by default. Single binary ships the
  embedded Vue 3 frontend — no separate Node server at runtime.
- **Frontend**: Vue 3 + Vite + antdv-next + UnoCSS + Pinia in `web` (a standalone
  package, no monorepo workspace). Build with `pnpm build` — Vite outputs directly into
  `internal/embedweb/dist`, which `go build` embeds.
- **CLI**: `cmd/openresty-admin/main.go` — cobra root, `cmd/openresty-admin/server.go` —
  HTTP server. The binary serves `/api/*` for the control plane and every other path from the
  embedded web dist with SPA fallback.
- **Go package layout**:
  - `internal/config/` — runtime configuration, defaults, env/flag loading.
  - `internal/nginx/` — the Host Adapter for OpenResty operations.
  - `internal/store/` — SQLite persistence matching the Node schema.
  - `internal/release/` — config change coordinator (snapshot → write → test → reload).
  - `internal/api/` — HTTP handlers, one file per domain + `server.go` wiring.
  - `internal/embedweb/` — embeds `web/dist` and serves it with SPA fallback.
- **Node dev mode**: `pnpm dev` still runs the Vue frontend with Vite dev server (HMR) and
  proxies `/api` to 3000. The Go binary can run alongside for API development.
- **Build**: `pnpm --dir web build` writes the frontend dist into `internal/embedweb/dist`, then `go build -o bin/openresty-admin ./cmd/openresty-admin/` embeds it. The resulting binary is self-contained.

## Commands

| Purpose | Command |
| --- | --- |
| Dev web only | `pnpm --dir web dev` |
| Build web dist | `pnpm --dir web build` |
| Build Go binary | `go build -o bin/openresty-admin ./cmd/openresty-admin/` |
| Build all | `pnpm --dir web build && go build -o bin/openresty-admin ./cmd/openresty-admin/` |
| Go tests | `go test ./internal/...` |
| Go vet | `go vet ./...` |

The `pnpm dev` and `pnpm typecheck` commands still work for the frontend. The Node
backend packages are no longer the primary server; they remain as reference only.

## Privileged operations — read before touching nginx

`internal/nginx/host.go` is the only Host Adapter for the OpenResty installation.
Treat every write path as a protected action:

- Route-level config changes must go through `release.Coordinator.Publish()` in
  `internal/release/coordinator.go`. It serializes releases in-process, snapshots the
  previous content, writes, runs `nginx -t`, reloads, and restores + re-tests on failure.
  Do not call `host.WriteConfig()` / `host.DeleteConfig()` from a route to persist a user edit.
- `Host.WriteConfig()` overwrites files under the nginx config root and snapshots the previous
  revision to `<file>.bak`; `Host.ReadForRelease()` reads with error propagation for the
  coordinator.
- `Host.Reload()` runs `<nginxBin> -s reload -c <nginxConfPath>` (`sudo -n` first, then
  unprivileged) with argv, never a shell string. `api/terminal.go` spawns a shell.
- `Host.DeleteConfig()` refuses directories and symlinks; it deletes single files only.

Required sequence for any config change:

1. State the intended change and the file it touches before writing.
2. `coordinator.Publish([]nginx.Change{{Path, Content, Delete?}}, "source")` — `Delete: true` removes.
3. If `nginx -t` fails the coordinator restores the previous content; never reload a failing config.
4. If reload fails the coordinator restores the previous content and re-tests it before reporting.
5. Report the resulting `config-releases` record. `recovery_required` means a human must intervene.

`routes/nginx.ts` returns 501 for stop/restart on purpose: process control has no release lock or
authentication yet. Keep it that way until both exist.

Stop conditions: never point the service at a production config root, never run these paths
against a host you were not asked to change, and never "fix" a failing `nginx -t` by reloading
anyway. Type checking catches none of this.

## Verification and debugging

- Pure logic (config parsing, storage, command assembly) belongs in `internal/*/` test files and
  runs via `go test ./internal/...`.
- Smallest loop for a Go change: `go vet ./... && go test ./internal/...`.
- Runtime probes: `GET /api/system/self-check` (nginx binary + `nginx -t`), `GET /api/system/info`,
  `GET /api/nginx/status`. The web `SelfCheck.vue` view renders the same surface.
  One-line liveness check: `curl -s http://127.0.0.1:3000/api/system/info`.
- Logs: the binary writes structured text to stdout; nginx access/error logs are exposed
  through `Host.ReadLog('access' | 'error', lines)` and `GET /api/logs/access` / `GET /api/logs/error`.

## Definition of done

The affected package passes `go vet`, the focused test covering the changed behaviour passes,
any nginx config change went through the sequence above, and new user-facing text goes through `t()`.

## Layout

Backend — Go
- `cmd/openresty-admin/main.go` — cobra CLI entry point, registers flags and `rootCmd`.
- `cmd/openresty-admin/server.go` — graceful HTTP listener with signal handling.
- `internal/config/config.go` — runtime configuration, default paths, env/flag loading.
- `internal/nginx/host.go` — the Host Adapter (openresty binary ops, config file I/O, logs).
- `internal/nginx/parser.go` — nginx config block parser (server, location + lua directives).
- `internal/nginx/backup.go` — config backup/restore via `.backups/` directory.
- `internal/nginx/helpers.go` — shared regexps and numeric helpers.
- `internal/store/store.go` — SQLite CRUD for collections, legacy JSON migration.
- `internal/release/coordinator.go` — config change coordinator with snapshot/validate/reload.
- `internal/api/server.go` — wires all route handlers, serves embedded frontend with SPA fallback.
- `internal/api/system.go` — dashboard, system info, nginx status/test/reload, self-check.
- `internal/api/sites.go` — site list/get/create/update/delete via config parsing.
- `internal/api/config.go` — config file CRUD, release history, backup history.
- `internal/api/collections.go` — generic CRUD for upstreams/streams/certificates/etc.
- `internal/api/seeds.go` — default seed data for collections.
- `internal/api/settings.go` — app settings with live path updates.
- `internal/api/logs.go` — log file list and content streaming.
- `internal/api/lua.go` — lua block enumeration and library listing.
- `internal/api/users.go` — user CRUD, profile, password change.
- `internal/api/notifications.go` — notifications list/create/mark-read/delete.
- `internal/api/backup.go` — backup list/create/restore/delete, auto-backup settings.
- `internal/api/terminal.go` — WebSocket terminal (spawns shell in nginx conf dir).
- `internal/api/sysinfo.go` — platform-specific system metrics (mem, cpu, uptime, load).
- `internal/embedweb/embed.go` — embeds `internal/embedweb/dist` (built by Vite) and serves it with SPA fallback.

Frontend — `web/src/`
- `routes/modules/<domain>.ts` — one file per route group, imported in `routes/index.ts`.
- `views/<domain>/` — Vue SFC with `<script setup lang="ts">`.
- `api/<domain>.ts` — typed wrappers over `api.get/post/put/delete`.
- `stores/` — `settings.ts` (theme, language, paths), `user.ts` (JWT auth), `app.ts`.
- `i18n/index.ts` — reactive `t(key)`, zh-CN/en-US.
- `composables/useMessage.ts`, `types/index.ts`, `types/route.d.ts`.

## Conventions

- antdv-next components are auto-imported in templates only; import `message`, `notification`, and
  `Modal` explicitly in `<script setup>`. Use `message.success()` / `message.error()`, never
  `AMessage.xxx()` in script.
- `dayjs` needs an explicit `import dayjs from 'dayjs'` — it is not auto-imported.
- Views fetch with `onMounted` + a `loading` ref + `ASpin`; destructive actions use `APopconfirm`.
- All user-facing text goes through `t()` from `@/i18n`.
- Backend CRUD goes through `loadCollection`, `saveCollection`, `createItem`, `updateItem`, and
  `deleteItem` from `services/storage.ts`. That module is SQLite-backed
  (`packages/server/data/openresty-admin.db` via `node:sqlite`; verified on Node 24, older Node
  22.x releases may need `--experimental-sqlite`); route code stays
  synchronous. Legacy `<name>.json` collections are imported once on first open, marked in the
  `meta` table, and left on disk as a backup — never re-imported after that. Override the data
  directory with `OPENRESTY_ADMIN_DATA_DIR`.
- `strict` is enabled in `web/tsconfig.json`.

## Frontend design language

`web` uses an **iOS 26 / Liquid Glass** visual language. It is centralized, not per-view:

- `src/styles/ios26.css` owns the tokens (`--bg`, `--glass`, `--group`, `--tint`, `--label-*`,
  the concentric radius scale), the material/typography primitives (`.ios-glass`, `.ios-group`,
  `.ios-row`, `.ios-title`, `.ios-instrument`, `.ios-chip`, `.ios-mono`), and the antdv surface
  overrides. **Change the language there, not in individual views.**
- `src/App.vue` bridges the same palette into antdv's own theme tokens.
- `src/layouts/DefaultLayout.vue` owns the shell: a floating glass rail, a sticky glass bar with a
  scroll-edge effect, and the content column. View headings marked `text-2xl` / `text-xl` are
  promoted to the large title by the shell.
- `src/components/EntryGroup.vue` is the shared grouped inset list used by the hub pages.

Rules: glass is only for chrome that floats (rail, bar, sheets, terminal) — content sits on opaque
inset groups with no drop shadow. Tint (`--tint`) means interactive; `--ok` / `--warn` / `--bad`
mean nginx state and are never used as decoration. Machine data (ports, paths, directives, log
lines) uses `.ios-mono`.

## Known gaps — do not assume these are handled

- The server has **no authentication middleware**: every `/api` route is reachable by anyone who
  can open the port, `routes/users.ts` compares plaintext passwords, and `routes/terminal.ts`
  spawns a shell. Keep the bind address on loopback until this is fixed.
- There is no CI configuration and no PR template; only the agent enforces the gate above.
- `packages/server/data/*` is untracked runtime state (the SQLite database, its WAL files, and
  legacy `users.json` etc.) — do not commit it.
