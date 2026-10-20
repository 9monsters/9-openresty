## ADR-003: Web UI Authentication — JWT

**Status:** Accepted  
**Date:** 2026-10-17

### Context

The management web UI must authenticate operators. Plaintext password comparison in the
current `users` collection is unacceptable for a multi-machine deployment where the hub
may listen on a non-loopback interface.

### Decision

Use **JWT bearer tokens** for web-to-hub authentication.

- Login endpoint `POST /api/auth/login` accepts `{username, password}`, returns
  `{token, expiresAt}`.
- All `/api/*` routes require `Authorization: Bearer <token>`.
- Tokens are HMAC-SHA256 signed (single-binary deployment, no need for asymmetric keys).
- Token expiry: 24 hours. A `POST /api/auth/refresh` endpoint extends the session.

For agent-to-hub mTLS, ADR-002 already covers it. This ADR addresses only the
browser-to-hub path.

### Implementation sketch

- `internal/auth/jwt.go` — sign/verify, extract claims.
- `internal/auth/middleware.go` — HTTP middleware, injects user identity into request context.
- `POST /api/auth/login` — bcrypt.CompareHashAndPassword against stored hash.
- `POST /api/auth/refresh` — accepts valid (non-expired) token, returns new token.
- `PUT /api/users/profile/password` — bcrypt.GenerateFromPassword on change.

### Migration

- On upgrade, existing plaintext passwords are bcrypt-hashed in-place the first time a user
  logs in (transparent migration).
- The `users` collection adds a `passwordHash` field; `password` becomes a deprecated
  plaintext field removed after migration.

### Consequences

- Every API request incurs a HMAC verification — negligible cost, no DB lookup.
- Token revocation is not immediate (no server-side blacklist). For admin-level operations
  like deleting users, a short token TTL (15 min) with refresh cycle mitigates risk.
- The embedded frontend's `api/client.ts` must attach the token to every request.