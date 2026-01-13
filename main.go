// ---------------- main.go ----------------
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go-sqlserver-api/config"
	"go-sqlserver-api/connections"
	"go-sqlserver-api/routes"

	"github.com/go-playground/validator/v10"
)

func main() {
	// ------------------ DATABASE ------------------
	db := connections.ConnectDB()
	defer db.Close()

	// ------------------ VALIDATOR ------------------
	validate := validator.New() // singleton

	// ------------------ ROUTER ------------------
	r := routes.RegisterRoutes(db, validate)

	// ------------------ SERVER ------------------
	srv := &http.Server{
		Addr:         config.GetPort(),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ------------------ GRACEFUL SHUTDOWN ------------------
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Server running on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server stopped gracefully")
}
