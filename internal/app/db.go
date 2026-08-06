package app

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/Alvesafk/shor/internal/models"
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

func (f FireDB) CreateURL(u models.URL, ctx context.Context) error {
	_, err := f.Collection("urls").Doc(u.ShortCode).Set(ctx, u)
	return err
}

func (f FireDB) DeleteURL(u models.URL, ctx context.Context) error {
	_, err := f.Collection("urls").Doc(u.ShortCode).Delete(ctx)
	return err
}

func (f FireDB) UpdateURL(n models.URL, oldShortCode string, ctx context.Context) error {
	_, err := f.Collection("urls").Doc(oldShortCode).Set(ctx, n, firestore.MergeAll)
	return err
}

func (f FireDB) GetURLByShortCode(shortCode string, ctx context.Context) (*models.URL, error) {
	var url models.URL
	doc, err := f.Collection("urls").Doc(shortCode).Get(ctx)
	if err != nil {
		return nil, err
	}

	err = doc.DataTo(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}
