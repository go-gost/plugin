package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedAuthenticatorServer
}

func (s *server) Authenticate(ctx context.Context, in *proto.AuthenticateRequest) (*proto.AuthenticateReply, error) {
	token := s.getCredentials(ctx)
	if token != "gost" {
		return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
	}

	reply := &proto.AuthenticateReply{}
	if in.GetUsername() == "gost" && in.GetPassword() == "gost" {
		reply.Ok = true
		reply.Id = "gost"
	}
	log.Printf("auth: %s %s, %s, %v", in.Client, in.GetUsername(), in.GetPassword(), reply.Ok)
	return reply, nil
}

func (s *server) getCredentials(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok && len(md["token"]) > 0 {
		return md["token"][0]
	}
	return ""
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterAuthenticatorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
