package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/limiter/traffic/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedLimiterServer
}

func (s *server) Limit(ctx context.Context, in *proto.LimitRequest) (*proto.LimitReply, error) {
	reply := &proto.LimitReply{
		In:  1024 * 1024,
		Out: 512 * 1024,
	}
	log.Printf("limiter: client=%s src=%s network=%s, addr=%s", in.Client, in.Src, in.Network, in.Addr)
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterLimiterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
