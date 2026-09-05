package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync/atomic"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

// The node types below mirror the JSON shape of x/config.NodeConfig. The hop
// reply is plain JSON on the wire, so the example stays a pure plugin-SDK
// consumer and doesn't import go-gost/x.
type authConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type connectorConfig struct {
	Type string       `json:"type"`
	Auth *authConfig `json:"auth,omitempty"`
}

type dialerConfig struct {
	Type string `json:"type"`
}

type nodeConfig struct {
	Name      string           `json:"name"`
	Addr      string           `json:"addr"`
	Connector *connectorConfig `json:"connector,omitempty"`
	Dialer    *dialerConfig    `json:"dialer,omitempty"`
}

var (
	nodes = []*nodeConfig{
		{
			Name: "node-0",
			Addr: ":8888",
			Connector: &connectorConfig{
				Type: "socks5",
				Auth: &authConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &dialerConfig{
				Type: "tcp",
			},
		},
		{
			Name: "node-1",
			Addr: ":9999",
			Connector: &connectorConfig{
				Type: "http",
				Auth: &authConfig{
					Username: "user",
					Password: "pass",
				},
			},
			Dialer: &dialerConfig{
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
	Src     string `json:"src"`
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
		log.Printf("hop: network=%s addr=%s host=%s client=%s src=%s", rb.Network, rb.Addr, rb.Host, rb.Client, rb.Src)

		node := nodes[counter.Add(1)%uint64(len(nodes))]
		json.NewEncoder(w).Encode(node)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
