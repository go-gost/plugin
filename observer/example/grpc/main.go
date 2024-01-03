package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-gost/plugin/observer/proto"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type server struct {
	proto.UnimplementedObserverServer
}

func (s *server) Observe(ctx context.Context, in *proto.ObserveRequest) (*proto.ObserveReply, error) {
	reply := &proto.ObserveReply{
		Ok: true,
	}

	for _, event := range in.Events {
		switch event.Type {
		case "service":
			log.Printf("observe %s: kind=%s, service=%s, status=%+v", event.Type, event.Kind, event.Service, event.Status)
		case "stats":
			log.Printf("observe %s: kind=%s, service=%s, client=%s, stats=%+v", event.Type, event.Kind, event.Service, event.Client, event.Stats)
		default:
			log.Printf("observe %+v", event)
		}
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
	proto.RegisterObserverServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
