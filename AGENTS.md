# Starter Redis implementation contract

This repository owns the independently versioned Redis integration for
Spice. Work directly on local `main` in bounded commits. Fetch before editing
and immediately before pushing; never overwrite unexpected remote work.

Go 1.26.5 is mandatory. Every product change must preserve secure connection
defaults, caller-owned contexts and pools, cache semantics, deterministic typed
serialization, expiration and pool isolation, and public core-contract
isolation. Add
positive and failure-path tests, update public documentation, run `make verify`
on the exact commit tree, and push only a green commit.

The normal gate is offline after dependencies are cached. Docker-backed
Redis acceptance is an additional release and hosted integration gate; it
may not replace deterministic unit and failure-path tests. Never commit live
database credentials or weaken TLS defaults to simplify production use.

Release-parity work must preserve the exact central release-tool versions
authorized by the root `go.mod`, invoke full package paths, and run both central
and retained rehearsals with workspace and network resolution disabled in vendor
mode. The retained repository builder and signed production workflow remain
authoritative until a separately reviewed signing migration; unsigned parity
must never manufacture signatures or key material.
