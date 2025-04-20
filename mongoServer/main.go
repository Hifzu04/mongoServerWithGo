// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/Hifzu04/myMongoServer/Router"
)

func main() {
	// 1. Load .env (optional)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: no .env (%v), proceeding with existing env vars", err)
	}

	// 2. Port from env (default 6000)
	port := os.Getenv("PORT")
	if port == "" {
		port = "6000"
	}

	// 3. Build handler (with CORS + routes)
	handler := router.Router()

	// 4. HTTP server with timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 5. Start in background
	go func() {
		log.Printf("🚀 Server listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// 6. Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("⚠️  Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}
	log.Println("✔️  Server stopped cleanly")
}
