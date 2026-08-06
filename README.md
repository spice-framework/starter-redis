# Spice Redis starter

`github.com/spice-framework/starter-redis` is the independently versioned,
opt-in Redis integration for Spice. It provides a secure, deterministic
go-redis client pool and a typed JSON implementation of Spice's cache
contract. Importing Spice core alone never selects go-redis or creates a
network client.

```go
client, cleanup, err := redis.Open(redis.Options{
    URL: "rediss://default:secret@redis.example.com:6379/0",
})
if err != nil {
    return err
}
defer cleanup(context.Background())

orders, err := redis.NewJSONStore[Order](client, redis.StoreOptions{
    Definition: cache.Definition{ID: "orders.by-id", Module: "orders"},
    Prefix:     "orders-by-id",
})
```

`Open` validates configuration and creates a caller-owned pool without network
I/O. `Ping` is the explicit readiness operation. `rediss` with TLS 1.2 or newer
and authenticated access are defaults; plaintext and unauthenticated local
development require separate opt-ins. URLs, credentials, keys, and encoded
values are never included in starter-generated failures or observations.

Each client owns an independent bounded pool and returns idempotent Spice
lifecycle cleanup. `JSONStore` provides deterministic namespacing, typed JSON,
bounded payloads, expiration, cancellation, local operation counters, and
module-aware observations without hidden global state.

## Install

```text
go get github.com/spice-framework/starter-redis@latest
```

During preview development, applications should pin an exact version covered
by [the compatibility contract](spice-compatibility.json).

## Verify

Go 1.26.5 is mandatory:

```text
make check
make release-parity
make verify
make verify-release
```

`make verify` checks the exact minimum and current Spice versions, then runs
the repository quality gate once: formatting, module/vendor reproducibility,
vet, allowlisted lint and nil safety, gosec, govulncheck, shuffled race tests,
at least 85% product coverage, and offline vendor builds. Compatibility tests
use isolated alternate modfiles and never rewrite committed dependency files.

Release parity runs the exact `spice-dev` tool authorized by `go.mod` and the
retained repository builder twice each, entirely from `vendor` with network and
workspace resolution disabled. It compares archive entries after normalizing
only the builders' documented root-directory spelling, requires equivalent
SBOM package and dependency facts, verifies both checksum files, and forbids
rehearsal signatures on Windows and Linux.

Real-system acceptance runs against the pinned Redis 8.4.0 multi-platform
image with password authentication. It covers independent pools, explicit
plaintext opt-in, ping, cancellation, typed JSON, expiration, deletion,
idempotent cleanup, and credential-safe authentication failure:

```text
docker run --detach --name spice-redis --publish 56379:6379 \
  redis@sha256:4eec4565e45aa0b3966554c866bc73211e281b0b3d89fe9a33c982e6faca809d \
  redis-server --requirepass spice-test --appendonly no
SPICE_REDIS_TEST_URL='redis://:spice-test@127.0.0.1:56379/0' make integration
docker rm --force spice-redis
```

See [the dependency review](docs/dependency-review.md) and
[support contract](docs/support.md) before production adoption.

## Releases

The repository builds deterministic source-only releases with an SPDX 2.3
SBOM, SHA-256 checksums, and Ed25519 signatures. See the exact artifact and
clean-tag ceremony in [the release guide](docs/releasing.md).
The retained repository builder and signed production workflow remain the
release authority while the centrally rendered unsigned candidate is held to
the dual-builder parity contract.
