package main

import (
	"log"
	"net/http"

	"github.com/BradiLugembaCS/payflow-api/internal/server"
)

func main() {
	srv := server.New()

	addr := ":8080"

	log.Printf("PayFlow API starting on %s", addr)

	err := http.ListenAndServe(addr, srv.Handler())
	if err != nil {
		log.Fatal(err)
	}
}
