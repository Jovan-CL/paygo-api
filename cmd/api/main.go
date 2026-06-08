package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"paygo-api/internal/config"
	"syscall"
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
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}

	// 4. Launch the Network Listener Inside an Asynchronous Goroutine Thread
	// This prevents ListenAndServe from blocking the application code execution flow below.
	go func() {
		fmt.Printf("Starting server on %s\n", serverAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 5. Establish the System Termination Signal Listen Channel
	shutdownChannel := make(chan os.Signal, 1)

	// Notify this channel if the user types Ctrl+C or the OS requests process closure
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGTERM)

	// This operation is a block read step. The program goes to sleep right here
	// until a terminal interrupt signal drops inside the communication channel.
	sig := <-shutdownChannel
	fmt.Printf("\n⚠️ Captured shutdown signal (%s). Launching cleanup routines...\n", sig)

	// 6. Enforce a Strict 10-Second Boundary for Existing Data Transactions
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	// Direct the server engine to stop taking new traffic and complete current actions
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error: Server forced to abort mid-execution loop: %v", err)
	} else {
		fmt.Println("🛑 HTTP Server successfully detached from network port.")
	}

	// 7. Securely Terminate Database Connections Natively
	// Since we are exiting the main function manually now, we explicitly call Close()
	// here to guarantee it finishes processing before the operating system process terminates.
	if err := db.Close(); err != nil {
		log.Printf("Error: Failed to safely close database connection pool: %v", err)
	} else {
		fmt.Println("🔌 PostgreSQL connection pool cleanly disconnected.")
	}

	fmt.Println("🏁 PayGo API Engine teardown complete. Goodbye!")

}
