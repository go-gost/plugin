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

type bypassRequest struct {
	Network string `json:"network"`
	Addr    string `json:"addr"`
	Host    string `json:"host"`
	Client  string `json:"client"`
}

type bypassResponse struct {
	OK bool `json:"ok"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/bypass", func(w http.ResponseWriter, r *http.Request) {
		rb := bypassRequest{}
		if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := bypassResponse{
			OK: true,
		}

		log.Printf("bypass: client=%s network=%s, addr=%s, host=%s", rb.Client, rb.Network, rb.Addr, rb.Host)

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
