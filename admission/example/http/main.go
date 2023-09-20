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

type admissionRequest struct {
	Addr string `json:"addr"`
}

type admissionResponse struct {
	OK bool `json:"ok"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/admission", func(w http.ResponseWriter, r *http.Request) {
		rb := admissionRequest{}
		if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := admissionResponse{}

		host := rb.Addr
		if v, _, _ := net.SplitHostPort(host); v != "" {
			host = v
		}
		if host == "127.0.0.1" {
			resp.OK = true
		}
		log.Printf("admission: %s, %v", rb.Addr, resp.OK)

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
