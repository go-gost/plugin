package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/go-gost/x/config"
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

type hopRequest struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Host    string `json:"host"`
	Client  string `json:"client"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	var counter atomic.Uint64

	http.HandleFunc("/hop", func(w http.ResponseWriter, r *http.Request) {
		rb := hopRequest{}
		if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("hop: network=%s addr=%s host=%s client=%s", rb.Network, rb.Addr, rb.Host, rb.Client)

		node := nodes[counter.Add(1)%uint64(len(nodes))]
		json.NewEncoder(w).Encode(node)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
