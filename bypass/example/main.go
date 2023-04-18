package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/bypass/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedBypassServer
}

func (s *server) Bypass(ctx context.Context, in *proto.BypassRequest) (*proto.BypassReply, error) {
	reply := &proto.BypassReply{}
	if in.GetAddr() == "example.com" {
		reply.Ok = true
	}
	log.Printf("bypass: %s, %v", in.GetAddr(), reply.Ok)
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterBypassServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
