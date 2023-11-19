package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/router/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedRouterServer
}

func (s *server) GetRoute(ctx context.Context, in *proto.GetRouteRequest) (*proto.GetRouteReply, error) {
	reply := &proto.GetRouteReply{}
	log.Printf("router get: dst=%s", in.GetDst())
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterRouterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
