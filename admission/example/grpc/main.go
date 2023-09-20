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

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedAdmissionServer
}

func (s *server) Admit(ctx context.Context, in *proto.AdmissionRequest) (*proto.AdmissionReply, error) {
	reply := &proto.AdmissionReply{}
	host := in.GetAddr()
	if v, _, _ := net.SplitHostPort(host); v != "" {
		host = v
	}
	if host == "127.0.0.1" {
		reply.Ok = true
	}
	log.Printf("admission: %s, %v", in.GetAddr(), reply.Ok)
	return reply, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAdmissionServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
