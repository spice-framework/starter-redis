module github.com/spice-framework/starter-redis

go 1.26.0

toolchain go1.26.5

require (
	github.com/redis/go-redis/v9 v9.21.0
	github.com/spice-framework/spice v0.1.0-preview.1
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/spice-framework/development v0.0.0-20260806132124-4c308d1b9fda // indirect
	github.com/spice-framework/toolchain v0.0.0-20260806133530-71211498297c // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/mod v0.38.0 // indirect
)

tool (
	github.com/spice-framework/development/cmd/spice-dev
	github.com/spice-framework/toolchain/cmd/spice-library-release-verify
)
