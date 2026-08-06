package app

import (
	"context"

	"cloud.google.com/go/firestore"
)

type FireDB struct {
	*firestore.Client
}

func DB(fireApp FireApp, ctx context.Context) (*FireDB, error) {
	fireDB, err := fireApp.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &FireDB{fireDB}, nil
}
