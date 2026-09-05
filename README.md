# GOST Plugin SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/go-gost/plugin.svg)](https://pkg.go.dev/github.com/go-gost/plugin)

gRPC-based plugin protocol definitions for [GOST](https://github.com/go-gost/gost). This module defines the service contracts that external plugin processes implement to extend GOST with custom authentication, admission control, traffic limiting, DNS resolution, and more.

## Available plugins

| Plugin | Service | Description |
|--------|---------|-------------|
| [admission](admission/) | `Admission` | Allow/deny connections by address |
| [auth](auth/) | `Authenticator` | Authenticate users and clients |
| [bypass](bypass/) | `Bypass` | Bypass the proxy chain for specific hosts |
| [hop](hop/) | `Hop` | Custom node selection for proxy chains |
| [hosts](hosts/) | `HostMapper` | Static or dynamic host-to-IP mapping |
| [ingress](ingress/) | `Ingress` | Ingress rule management |
| [limiter/traffic](limiter/traffic/) | `Limiter` | Per-connection bandwidth limits |
| [observer](observer/) | `Observer` | Receive connection and traffic events |
| [p2p](p2p/) | `P2P` | Establish tunnels to peers for P2P connectivity |
| [recorder](recorder/) | `Recorder` | Record proxied traffic |
| [resolver](resolver/) | `Resolver` | Custom DNS resolution |
| [router](router/) | `Router` | Custom route selection |
| [sd](sd/) | `SD` | Service discovery (register/deregister/get) |

## Usage

### gRPC plugin server

```go
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "net"

    "github.com/go-gost/plugin/admission/proto"
    "google.golang.org/grpc"
)

var port = flag.Int("port", 8000, "Server port")

type server struct {
    proto.UnimplementedAdmissionServer
}

func (s *server) Admit(ctx context.Context, req *proto.AdmissionRequest) (*proto.AdmissionReply, error) {
    reply := &proto.AdmissionReply{}
    if req.GetAddr() == "127.0.0.1:80" {
        reply.Ok = true
    }
    return reply, nil
}

func main() {
    flag.Parse()
    lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    s := grpc.NewServer()
    proto.RegisterAdmissionServer(s, &server{})
    log.Fatal(s.Serve(lis))
}
```

### HTTP plugin server

GOST also supports HTTP/JSON transport. The endpoint path matches the plugin type:

```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net"
    "net/http"
)

var port = flag.Int("port", 8000, "Server port")

func main() {
    flag.Parse()
    lis, _ := net.Listen("tcp", fmt.Sprintf(":%d", *port))
    http.HandleFunc("/admission", func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Addr    string `json:"addr"`
            Service string `json:"service"`
        }
        json.NewDecoder(r.Body).Decode(&req)
        json.NewEncoder(w).Encode(map[string]bool{"ok": req.Addr == "127.0.0.1:80"})
    })
    log.Fatal(http.Serve(lis, nil))
}
```

### Configuration

Point GOST at your plugin server via the `plugin` block in config:

```yaml
services:
  - name: service-0
    addr: ":8080"
    handler:
      type: http
    listener:
      type: tcp
    plugin:
      type: admission
      addr: 127.0.0.1:8000
```

### P2P plugin

The `P2P` service lets GOST establish the network path to a chain node through
a tunnel opened by an external plugin (NAT traversal, relay, hole punching —
the strategy is entirely up to the plugin). The client sends an opaque `peer`
identifier; the plugin replies with an opaque `endpoint` handle (v1: a locally
dialable TCP `host:port`). P2P plugins are declared in a top-level `p2ps:`
section and referenced by chain node metadata:

```yaml
p2ps:
  - name: p2p-1
    plugin:
      type: grpc
      addr: 127.0.0.1:8003

chains:
  - name: chain-0
    hops:
      - name: hop-0
        nodes:
          - name: node-0
            addr: 192.168.1.10:8080
            connector:
              type: http
            dialer:
              type: tcp
            metadata:
              p2p: p2p-1
```

A reference implementation of the plugin host (stub: local TCP bridge) lives in
the standalone [p2p](https://github.com/go-gost/p2p) repo.

## Code generation

Generated protobuf code was produced with `protoc-gen-go v1.28.1` and `protoc-gen-go-grpc v1.2.0`. To regenerate:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    admission.proto
```

## License

MIT
