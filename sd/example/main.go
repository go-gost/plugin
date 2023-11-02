package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/sd/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedSDServer
}

func (s *server) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterReply, error) {
	reply := &proto.RegisterReply{}
	log.Printf("register: %+v", in.GetService())
	return reply, nil
}

func (s *server) Deregister(ctx context.Context, in *proto.DeregisterRequest) (*proto.DeregisterReply, error) {
	reply := &proto.DeregisterReply{}
	log.Printf("deregister: %+v", in.GetService())
	return reply, nil
}

func (s *server) Renew(ctx context.Context, in *proto.RenewRequest) (*proto.RenewReply, error) {
	reply := &proto.RenewReply{}
	log.Printf("renew: %+v", in.GetService())
	return reply, nil
}

func (s *server) Get(ctx context.Context, in *proto.GetServiceRequest) (*proto.GetServiceReply, error) {
	log.Printf("get: %s", in.GetName())
	reply := &proto.GetServiceReply{
		Services: []*proto.Service{
			{
				Name:    "c25cab37-4d84-4c34-8c30-4916a26b8fe9",
				Network: "tcp",
				Address: "127.0.0.1:8420",
			},
			{
				Name:    "b83294f5-34b6-4698-b15d-b3c11b81ec1c",
				Network: "udp",
				Address: "127.0.0.1:8420",
			},
		},
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
	proto.RegisterSDServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
