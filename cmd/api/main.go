package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/BradiLugembaCS/payflow-api/internal/database"
	"github.com/BradiLugembaCS/payflow-api/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL")

	srv := server.New()

	addr := "127.0.0.1:8080"

	log.Printf("PayFlow API starting on %s", addr)

	if err := http.ListenAndServe(addr, srv.Handler()); err != nil {
		log.Fatal(err)
	}
}
