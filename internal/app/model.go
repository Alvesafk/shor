package app

import (
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
