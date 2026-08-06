package app

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type URL struct {
	ID          uuid.UUID `json:"id"`
	URL         string    `json:"url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccessCount uint      `json:"access_count"`
}

func GenURLStruct(url, shortCode string) *URL {
	defer func() {
		if err := recover(); err != nil { // Need this recover because uuid.New() can panic.
			log.Println("Recovered from panic: ", err)
		}
	}()

	id := uuid.New()

	r := URL{
		ID:          id,
		URL:         url,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AccessCount: 0,
	}

	return &r
}
