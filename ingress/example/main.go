package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/ingress/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedIngressServer
}

func (s *server) Get(ctx context.Context, in *proto.GetRequest) (*proto.GetReply, error) {
	reply := &proto.GetReply{}
	log.Printf("ingress get: %s", in.GetHost())
	return reply, nil
}

func (s *server) Set(ctx context.Context, in *proto.SetRequest) (*proto.SetReply, error) {
	reply := &proto.SetReply{}
	log.Printf("ingress set: %s -> %s", in.GetHost(), in.GetEndpoint())
	return reply, nil

}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterIngressServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
