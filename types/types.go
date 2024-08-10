package types

import (
	"time"
)

type Account struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"createdAt"`
	ShortURL    string    `json:"shortUrl"`
	OriginalURL string    `json:"originalUrl"`
}

func GenerateNewAccount(name, email, originalURL string) *Account {
	return &Account{
		Name:        name,
		Email:       email,
		CreatedAt:   time.Now().UTC(),
		OriginalURL: originalURL,
	}
}
