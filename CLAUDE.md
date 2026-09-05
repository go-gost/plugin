# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Code Generation

```bash
# Build the module
go build ./...

# Regenerate protobuf code (run from a service proto directory)
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    admission.proto
```

Generated code was produced with `protoc-gen-go v1.28.1`, `protoc-gen-go-grpc v1.2.0`, `protoc v3.15.8`. The `go_package` option in each `.proto` sets the output path (e.g. `github.com/go-gost/plugin/admission/proto`).

## Architecture

This module defines the **gRPC plugin protocol** for [GOST](https://github.com/go-gost/gost). External plugin processes implement these gRPC services to extend GOST's behavior (authentication, rate limiting, DNS resolution, etc.). The core GOST binary communicates with plugins either via **gRPC** or **HTTP/JSON** — examples for both transports are provided.

### Service catalog (13 plugin types)

Each subdirectory follows the same structure: `proto/<name>.proto` + generated `.pb.go`/`_grpc.pb.go`, and `example/` with reference implementations (gRPC and optionally HTTP).

| Directory | Service | RPC(s) | In config `plugin:` block |
|---|---|---|---|
| `admission/` | `Admission` | `Admit(req) → (ok)` | `type: admission` |
| `auth/` | `Authenticator` | `Authenticate(req) → (ok, id)` | `type: auth` |
| `bypass/` | `Bypass` | `Bypass(req) → (ok)` | `type: bypass` |
| `hop/` | `Hop` | `Select(req) → (node bytes)` | `type: hop` |
| `hosts/` | `HostMapper` | `Lookup(req) → (ips[], ok)` | `type: hosts` |
| `ingress/` | `Ingress` | `GetRule(req) → (endpoint)`, `SetRule(req) → (ok)` | `type: ingress` |
| `limiter/traffic/` | `Limiter` | `Limit(req) → (in, out)` | `type: limiter` |
| `observer/` | `Observer` | `Observe(events[]) → (ok)` | `type: observer` |
| `p2p/` | `P2P` | `OpenTunnel(req) → (ok, id, endpoint)`, `CloseTunnel(id) → (ok)`, `Status` | top-level `p2ps:` section (`plugin.type: grpc`) |
| `recorder/` | `Recorder` | `Record(data, metadata) → (ok)` | `type: recorder` |
| `resolver/` | `Resolver` | `Resolve(req) → (ips[], ok)` | `type: resolver` |
| `router/` | `Router` | `GetRoute(req) → (dst, gateway)` | `type: router` |
| `sd/` | `SD` | `Register`, `Deregister`, `Renew`, `Get` | `type: sd` |

All RPCs are **unary** (no streaming). Every request carries proxy context fields (`network`, `addr`, `host`, `client`, `service`) relevant to the intercepted connection. The p2p service deviates deliberately: `peer` and `endpoint` are opaque strings (the plugin defines their semantics), there is no `example/` in this module (the standalone `p2p/` host repo is the reference implementation), and only the gRPC transport exists.

### Dual transport: gRPC + HTTP

GOST supports two transports for plugin communication. Examples for both exist under `example/grpc/` and `example/http/`.

- **gRPC examples**: Standard gRPC server embedding the `Unimplemented*Server` struct. Auth example shows how to extract credentials from gRPC metadata.
- **HTTP examples**: Plain HTTP/JSON servers. The endpoint path matches the service name (e.g., `/admission`, `/auth`, `/bypass`). Request/response structs mirror the proto messages as JSON.

### Adding a new plugin type

The pattern is consistent across all 12 services:

1. Create `<name>/proto/<name>.proto` with the service, RPCs, and messages. Set `go_package` to `github.com/go-gost/plugin/<name>/proto`.
2. Run `protoc` to generate `.pb.go` and `_grpc.pb.go`.
3. Add at minimum a gRPC example under `<name>/example/grpc/main.go`.
4. If the plugin type supports HTTP transport, add `<name>/example/http/main.go`.

The corresponding implementation in `x/` (e.g., `x/admission/`) contains a gRPC client that calls these services and an HTTP client that hits the JSON endpoints — the dual transport logic lives there, not in this module.

### Known issue: recorder go_package

`recorder/proto/recorder.proto` declares `go_package = "github.com/go-gost/plugin/ingress/proto"`, which places it in the same Go package as the ingress service. This means `recorder/` and `ingress/` share a proto package namespace. If you add new message types to either, ensure no symbol collisions occur.

## Verification

```bash
# Build check from plugin/ directory
go build ./...

# Vet
go vet ./...
```

There are no tests in this module. Building is the primary verification.
