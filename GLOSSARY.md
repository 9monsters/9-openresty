# OpenResty Admin — Glossary

<!-- Last updated: 2026-10-20 -->
<!-- This glossary is canonical. When a term appears in code, docs, or UI, it MUST match the definition below. -->

## Core Entities

### Site
A single nginx `server { }` block — identified by `server_name` and `listen` directives, may contain Location children, SSL certificate bindings, and Lua lifecycle hooks. Equivalent to "站点" in the UI.

### Location
A single nginx `location <path> { }` block nested inside a Site. Defines request routing (`proxy_pass`, `root`, `return`) and optional Lua/WAF hooks scoped to a URL path prefix.

### Upstream
A named `upstream { }` block defining a group of backend servers with load‑balancing strategy (round‑robin, least‑conn, ip‑hash). Referenced by Site or Location `proxy_pass` directives.

### Stream
A TCP or UDP `server { }` block inside the nginx `stream { }` context — Layer‑4 proxy. Has listen address, proxy pass target, and optional SSL.

### Certificate
An SSL/TLS key‑pair (PEM) identified by domain names. Can be manually imported or managed via ACME. Bound to a Site’s `listen 443 ssl` directive.

### Node (DEPRECATED — single‑machine only)
~~A managed OpenResty instance identified by hostname/IP.~~ Removed in the single‑machine pivot. The local host is the only managed target; no remote Agent communication.

---

## Release & Change Management

### Release
An atomic configuration publish operation. A Release moves through the state machine
`Draft → Preparing → Validating → Reloading → HealthChecking → Succeeded`
and can transition to `RolledBack` or `RecoveryRequired` on failure.
One Release contains one or more Changes.

### Change
A single file‑level mutation within a Release: **create**, **update**, or **delete** of one configuration file under the nginx config root. Each Change carries a path, new content (empty for delete), and optional checksum.

### Draft
A pending configuration edit that has been saved but NOT published. Drafts are private to the author and do not affect the running nginx configuration. Publishing a Draft creates a Release.

### RecoveryRequired
A terminal failure state of a Release: automatic rollback failed, the configuration files on disk may be inconsistent, and the Release coordinator is globally locked. Only manual intervention (human confirmation + passing `nginx -t`) can clear this state and unlock new releases.

---

## WAF (Web Application Firewall)

### WAF Rule
A single security check executed inside `access_by_lua*`. Defines a **match condition** (IP/CIDR, header, body, URI, query‑arg pattern) and an **action** (deny, log, redirect, allow, challenge). Rules are Lua functions generated from structured parameters.

### WAF Policy
An ordered collection of WAF Rules with a shared **mode** (`off` | `observe` | `block`) and **scope** (global, per‑Site, per‑Location). Policies are bound to one or more Sites/Locations and carry a priority for conflict resolution.

### WAF Challenge
A WAF Rule action that requires the client to prove it is not a bot — typically a JavaScript computational challenge or CAPTCHA. The challenge is injected via `access_by_lua*` and the client is redirected or blocked until it passes.

---

## Lua Management

### Lua Script
A versioned unit of Lua code authored by the operator. Has a **name**, **semantic version**, **lifecycle stage** (`access_by_lua`, `log_by_lua`, etc.), and **scope** (`global` | `server` | `location`). Stored as a database entity; rendered into nginx `*_by_lua_file` directives by the release pipeline.

### Lua Stage
The OpenResty lifecycle phase where a Lua Script is mounted. Examples: `init_by_lua`, `rewrite_by_lua`, `access_by_lua`, `content_by_lua`, `log_by_lua`. The stage determines execution order and available nginx API capabilities.

### Lua Scope
The nginx context where a Lua Script executes: **global** (http block), **server** (server block), or **location** (location block). Scope constrains which requests the script applies to.

---

## DNS

### DNS Credential
API credentials for a cloud DNS provider (Alibaba Cloud DNS, Tencent Cloud DNSPod). Stored encrypted at rest. Used by DDNS to update records and by the certificate ACME DNS‑01 challenge.

### DNS Domain
A registered domain zone managed through a DNS provider. Contains DNS records (A, AAAA, CNAME, MX, TXT, etc.). The tool reads/updates records via the provider API.

### DDNS (Dynamic DNS)
A DNS record whose IP address is periodically updated — typically used for home or dynamic‑IP endpoints. The tool polls the host’s public IP and pushes updates to the configured DNS provider.

---

## RBAC

### Role
One of three permission levels: **Viewer** (read‑only dashboard, logs, status), **Operator** (Viewer + configuration editing, draft creation, release publishing), **Administrator** (Operator + user management, process control, WAF policy enable/disable, settings changes).

### Permission
Not a standalone database entity. Permissions are derived from Role at request time by the auth middleware. Future: if finer granularity is needed, Roles can be decomposed into individual capabilities.

---

## Infrastructure

### Host Adapter
The ONLY package allowed to touch the OpenResty installation: config files, the nginx binary, log files, and PID file. All path resolution rejects symlinks, `..` escapes, and non‑regular files. All commands use argument arrays, never shell strings.

### Release Coordinator
The singleton service that serialises all configuration changes. It owns the Release state machine, enforces the mutual‑exclusion lock, and is the sole caller of Host Adapter write/test/reload operations. No route handler writes config files directly.

### Health Probe
An active HTTP request sent to each affected `server_name` after `nginx -s reload` succeeds. The probe expects a response with HTTP status < 500. If any probe fails, the Release coordinator initiates rollback.
