package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("Connecting to Firebase.")
	fireApp, err := app.New()
	if err != nil {
		log.Fatal("error: ", err)
	}
	fmt.Println("Connected.")

	fireDB, err := app.DB(*fireApp, ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}

	conn := handlers.NewConnection(fireDB, ctx)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", conn.GetHelloWorld)
	r.Post("/shor", conn.PostURL)

	fmt.Println("Server listening on port", serverPort)
	http.ListenAndServe(serverPort, r)
}
