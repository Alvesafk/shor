package main

import (
	"context"
	"log"

	"github.com/Alvesafk/shor/internal/db"
)

func main() {
	ctx := context.Background()

	fireApp, err := db.Connect()
	if err != nil {
		log.Fatal("error: ", err)
	}

	fireDB, err := fireApp.Firestore(ctx)
	if err != nil {
		log.Fatal("error: ", err)
	}

	_, err = fireDB.Collection("hello").Doc("world").Set(ctx, map[string]any{
		"hello":  "World",
		"world": "Hello",
	})
}
