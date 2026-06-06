package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"paygo-api/internal/config"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Printf("Connected to database successfully\n")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{
			"status":    "ok",
			"message":   "PayGo API is healthy",
			"version":   "1.0.0",
			"timestamp": fmt.Sprintf("%d", time.Now().Unix()),
		}
		json.NewEncoder(w).Encode(response)
	})

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Starting server on %s\n", serverAddr)

	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
