package main

import (
	"context"
	"log"

	"github.com/Alvesafk/shor/internal/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"
)

func main() {
	ctx := context.Background()

	fireApp, err := app.New()
	if err != nil {
		log.Fatal("error: ", err)
	}

	fireDB, err := app.DB(*fireApp, ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, err = fireDB.Collection("hello").Doc("world").Set(ctx, map[string]any{
		"hello":  "World",
		"world": "Hello",
	})

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
}
