package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type limiterRequest struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Client  string `json:"client"`
	Src     string `json:"src"`
}

type limiterResponse struct {
	In  int64 `json:"in"`
	Out int64 `json:"out"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/limiter", func(w http.ResponseWriter, r *http.Request) {
		req := limiterRequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := limiterResponse{
			In:  1024 * 1024, // 1MB
			Out: 512 * 1024,  // 512KB
		}

		log.Printf("limiter: client=%s src=%s network=%s, addr=%s", req.Client, req.Src, req.Network, req.Addr)

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
