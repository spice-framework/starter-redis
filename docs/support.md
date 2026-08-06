# Support and compatibility

| Contract | Current development support |
|---|---|
| Go | Exactly 1.26.5 for development and release verification |
| Spice minimum/current | Exact versions in [`spice-compatibility.json`](../spice-compatibility.json) |
| Redis | 8.4.0 real-system acceptance; server versions supported by go-redis v9.21.0 remain integration targets, not yet release claims |
| Operating systems | Windows, Linux, and macOS; Linux container acceptance |
| Architectures | amd64 and arm64 compilation through the public core API |
| Transport security | `rediss` with TLS 1.2 minimum; plaintext requires explicit opt-in |
| Authentication | Required by default; unauthenticated access requires a separate explicit opt-in |
| Real-system artifact | `redis:8.4.0-alpine3.22` index digest `sha256:4eec4565e45aa0b3966554c866bc73211e281b0b3d89fe9a33c982e6faca809d` |
| Central release tools | `github.com/spice-framework/development/cmd/spice-dev` at `v0.0.0-20260806121906-963bb6676069`; `github.com/spice-framework/toolchain/cmd/spice-library-release-verify` at `v0.0.0-20260806054457-a83d9b58034c` |

The first preview tag will define the first published minimum Spice version.
Until then, `spice-compatibility.json` is the sole compatibility boundary
source. Its minimum must equal the exact direct Spice requirement in `go.mod`;
its current value is a forward-compatibility endpoint, not an unbounded runtime
dependency. The compatibility runner uses isolated alternate module files and
does not mutate the committed module or vendor graph.

A release may raise the minimum only through an intentional `go.mod` change,
an updated compatibility contract, and green minimum/current jobs. A moving
branch name is never written to release metadata or used by applications.

The pinned central tool renders unsigned rehearsal candidates only. Windows
and Linux CI compare them with the retained builder under vendor-only offline
resolution; the retained command remains the signed production authority.
