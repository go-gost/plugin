package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

var (
	port = flag.Int("port", 8000, "The server port")
)

type observeRequest struct {
	Events []event `json:"events"`
}

type event struct {
	Kind    string       `json:"kind"`
	Service string       `json:"service"`
	Client  string       `json:"client,omitempty"`
	Type    string       `json:"type"`
	Stats   *statsEvent  `json:"stats,omitempty"`
	Status  *statusEvent `json:"status,omitempty"`
}

type statsEvent struct {
	TotalConns   uint64 `json:"totalConns"`
	CurrentConns uint64 `json:"currentConns"`
	InputBytes   uint64 `json:"inputBytes"`
	OutputBytes  uint64 `json:"outputBytes"`
}

type statusEvent struct {
	State string `json:"state"`
	Msg   string `json:"msg"`
}

type observeResponse struct {
	OK bool `json:"ok"`
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v", lis.Addr())

	http.HandleFunc("/observer", func(w http.ResponseWriter, r *http.Request) {
		rb := observeRequest{}
		data, _ := io.ReadAll(r.Body)
		log.Println(string(data))
		if err := json.NewDecoder(bytes.NewReader(data)).Decode(&rb); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := observeResponse{
			OK: true,
		}

		for _, event := range rb.Events {
			switch event.Type {
			case "service":
				log.Printf("observe %s: kind=%s, service=%s, status=%+v", event.Type, event.Kind, event.Service, event.Status)
			case "stats":
				log.Printf("observe %s: kind=%s, service=%s, client=%s, stats=%+v", event.Type, event.Kind, event.Service, event.Client, event.Stats)
			default:
				log.Printf("observe %+v", event)
			}
		}

		json.NewEncoder(w).Encode(resp)
	})

	if err := http.Serve(lis, nil); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
