package app

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/Alvesafk/shor/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UrlNotFound = errors.New("not found")
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
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("not found")
		}

		return nil, err
	}

	err = doc.DataTo(&url)
	if err != nil {
		return nil, err
	}
	url.ID = doc.Ref.ID

	return &url, nil
}

func (f FireDB) ShortURLExists(shortCode string, ctx context.Context) (bool, error) {
	_, err := f.GetURLByShortCode(shortCode, ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return false, UrlNotFound
		}

		return false, err
	}

	return true, nil
}

func (f FireDB) URLAlreadyRegistered(urlString string, ctx context.Context) (*models.URL, bool, error) {
	docs, err := f.Collection("urls").
		Where("URL", "==", urlString).
		Limit(1).
		Documents(ctx).GetAll()
	if err != nil {
		return nil, false, err
	}

	if len(docs) > 0 {
		var url models.URL
		if err := docs[0].DataTo(&url); err != nil {
			return nil, false, err
		}

		url.ID = docs[0].Ref.ID
		return &url, true, nil
	}

	return nil, false, nil
}
