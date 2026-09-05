package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/go-gost/plugin/hop/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

// The node types below mirror the JSON shape of x/config.NodeConfig. The
// SelectReply node is plain JSON bytes on the wire, so the example stays a
// pure plugin-SDK consumer and doesn't import go-gost/x.
type authConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type connectorConfig struct {
	Type string       `json:"type"`
	Auth *authConfig `json:"auth,omitempty"`
}

type dialerConfig struct {
	Type string `json:"type"`
}

type nodeConfig struct {
	Name      string           `json:"name"`
	Addr      string           `json:"addr"`
	Connector *connectorConfig `json:"connector,omitempty"`
	Dialer    *dialerConfig    `json:"dialer,omitempty"`
}

var (
	nodes = []*nodeConfig{
		{
			Name: "node-0",
			Addr: ":8888",
			Connector: &connectorConfig{
				Type: "socks5",
				Auth: &authConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &dialerConfig{
				Type: "tcp",
			},
		},
		{
			Name: "node-1",
			Addr: ":9999",
			Connector: &connectorConfig{
				Type: "http",
				Auth: &authConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &dialerConfig{
				Type: "tcp",
			},
		},
	}
)

type server struct {
	counter atomic.Uint64
	proto.UnimplementedHopServer
}

func (s *server) Select(ctx context.Context, in *proto.SelectRequest) (*proto.SelectReply, error) {
	node := nodes[s.counter.Add(1)%uint64(len(nodes))]
	v, _ := json.Marshal(node)
	reply := &proto.SelectReply{
		Node: v,
	}
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterHopServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
