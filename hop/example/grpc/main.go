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
	"github.com/go-gost/x/config"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

var (
	nodes = []*config.NodeConfig{
		{
			Name: "node-0",
			Addr: ":8888",
			Connector: &config.ConnectorConfig{
				Type: "socks5",
				Auth: &config.AuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &config.DialerConfig{
				Type: "tcp",
			},
		},
		{
			Name: "node-1",
			Addr: ":9999",
			Connector: &config.ConnectorConfig{
				Type: "http",
				Auth: &config.AuthConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &config.DialerConfig{
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
