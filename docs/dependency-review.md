# Dependency review: go-redis

- Decision: approved for the independently versioned
  `github.com/spice-framework/starter-redis` module.
- Version: `github.com/redis/go-redis/v9` v9.21.0.
- Upstream: <https://github.com/redis/go-redis>.
- License: BSD-2-Clause; retained in the vendored module license.
- Maintenance: the official Redis Go client has an active stable v9 line and
  supports the current Redis releases.
- Security: complete URLs prevent environment fallback, `rediss` with TLS 1.2
  or newer is the default, authentication is required by default, local
  plaintext and unauthenticated use require separate explicit opt-ins, query
  options are rejected, and failures never include the URL or password.
  `gosec` and `govulncheck` cover the reachable dependency graph.
- Cancellation: every network operation uses a caller-owned context with
  context deadlines enabled in the client.
- Observability: the client name is bounded and validated for safe server
  attribution. The typed cache adapter emits module-aware Spice observations
  without keys or values; optional OpenTelemetry hooks remain
  application-owned.
- Configuration: Spice fixes protocol, retry, timeout, and bounded-pool
  defaults instead of inheriting CPU- or process-environment-dependent
  behavior. `Open` performs no network I/O and returns exact lifecycle cleanup.
- Integration: a tagged test runs against the official Redis 8.4.0 Alpine
  multi-platform index (`redis:8.4.0-alpine3.22`,
  `sha256:4eec4565e45aa0b3966554c866bc73211e281b0b3d89fe9a33c982e6faca809d`)
  and exercises authenticated connection, independent pools, ping, bounded
  typed JSON cache operations, expiration, cancellation, idempotent cleanup,
  and credential-safe authentication failures on amd64 CI. The reviewed index
  also provides arm64.

Primary references:

- <https://github.com/redis/go-redis/releases/tag/v9.21.0>
- <https://redis.io/docs/latest/develop/clients/go/connect/>
- <https://redis.io/docs/latest/develop/clients/go/error-handling/>

## Build-only dependencies: central release tools

- Decision: approved only as repository-authorized release tooling.
- Renderer: `github.com/spice-framework/development/cmd/spice-dev` from
  `github.com/spice-framework/development`
  `v0.0.0-20260806132124-4c308d1b9fda`.
- Independent verifier:
  `github.com/spice-framework/toolchain/cmd/spice-library-release-verify` from
  `github.com/spice-framework/toolchain`
  `v0.0.0-20260806133530-71211498297c`.
- Tool registration: both commands use standard Go `tool` directives and all
  invocations use their full package paths.
- License: Apache-2.0, with its notice retained in `vendor`.
- Runtime scope: none. Product packages import neither tool module, and
  released applications acquire no runtime dependency on them.
- Dependency graph: both tools participate in normal Go minimal-version
  selection. That build-time coupling is accepted and visible in `go.mod`,
  `go.sum`, and `vendor/modules.txt`; no parallel tool registry is introduced.
- Integrity and network behavior: both exact pseudo-versions are pinned and
  checksummed. Release parity runs with `GOWORK=off`, `GOPROXY=off`,
  `GOTOOLCHAIN=local`, and `GOFLAGS=-mod=vendor`, so it cannot select an ambient
  checkout, upgrade itself, or download dependencies.
- Security: the trusted native renderer reads the exact committed Git graph
  and writes only to caller-supplied temporary output directories. The
  verifier independently checks release artifacts without signing them. The
  rehearsal emits no signatures or signing material.
- Maintenance: the pinned central workflow owns production signing. The retained
  local builder stays only as a parity oracle until separate removal evidence is
  reviewed.
