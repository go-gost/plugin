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

type autherRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type autherResponse struct {
	OK bool   `json:"ok"`
	ID string `json:"id"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		rb := autherRequest{}
		if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := autherResponse{}
		if rb.Username == "gost" && rb.Password == "gost" {
			resp.OK = true
			resp.ID = "gost"
		}

		log.Printf("auth: %s %s, %s, %v", rb.Client, rb.Username, rb.Password, resp.OK)

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
