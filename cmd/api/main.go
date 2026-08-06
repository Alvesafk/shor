package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alvesafk/shor/internal/app"
	"github.com/Alvesafk/shor/internal/handlers"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	serverPort = ":7000"
)

func main() {
	ctx := context.Background()

	log.Println("Connecting to Firebase.")
	fireApp, err := app.New()
	if err != nil {
		log.Fatalf("Error on connecting to Firebase: %v", err)
	}
	log.Println("Connected.")

	log.Println("Connecting to Firestore.")
	fireDB, err := app.DB(*fireApp, ctx)
	if err != nil {
		log.Fatalf("Error on connecting to Firestore: %v", err)
	}
	log.Println("Connected.")

	conn := handlers.NewConnection(fireDB, ctx)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", conn.GetHelloWorld)

	r.Post("/shor", conn.PostURL)
	r.Get("/shor/{shortUrl}", conn.GetURL)

	server := &http.Server{
		Addr:         serverPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("Server is running on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("Server error: %v", err)

	case sig := <-stop:
		fmt.Println()
		log.Printf("Received shutdown signal: %v", sig)

	}

	log.Println("Server is shutting down...")

	ctxShutdown, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server shutdown error: %v", err)
		return
	}

	log.Println("Server exited properly.")
}
