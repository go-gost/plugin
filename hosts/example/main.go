package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/hosts/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedHostMapperServer
}

func (s *server) Lookup(ctx context.Context, in *proto.LookupRequest) (*proto.LookupReply, error) {
	reply := &proto.LookupReply{}
	if in.GetHost() == "localhost" {
		reply.Ips = []string{"127.0.0.1"}
		reply.Ok = true
	}
	log.Printf("hosts: %s/%s, %v", in.GetHost(), in.GetNetwork(), reply.Ok)
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterHostMapperServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
