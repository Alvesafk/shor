package db

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type FireApp struct {
	*firebase.App
}

func Connect() (*FireApp, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	serviceKey := os.Getenv("SERVICE_ACCOUNT_KEY")
	if serviceKey == "" {
		return nil, fmt.Errorf("could not find service account json file")
	}

	ctx := context.Background()

	opt := option.WithAuthCredentialsFile(option.ServiceAccount, serviceKey)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return &FireApp{app}, nil
}
