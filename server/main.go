package main

import (
	"fmt"
	"log"
	"net/http"

	v1 "github.com/seitamuro/go-auth0-2/server/handlers/v1"
)

const (
	port = 8000
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1", v1.HandleIndex)

	addr := fmt.Sprintf(":%d", port)

	log.Printf("Listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
