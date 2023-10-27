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

type recorderRequest struct {
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

type recorderResponse struct {
	OK bool `json:"ok"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/recorder", func(w http.ResponseWriter, r *http.Request) {
		rb := recorderRequest{}
		if err := json.NewDecoder(r.Body).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := recorderResponse{
			OK: true,
		}

		log.Printf("recorder: data=%s, metadata: %s", string(rb.Data), string(rb.Metadata))

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
