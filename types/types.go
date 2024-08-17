package types

import (
	"time"
)

type Account struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	ShortURL     string    `json:"shortUrl"`
	OriginalURL  string    `json:"originalUrl"`
}

type ShortURL struct {
	ID          int       `json:"id"`
	AccountID   int       `json:"accountId"`
	OriginalURL string    `json:"originalUrl"`
	ShortURL    string    `json:"shortUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}

func GenerateNewAccount(name, email, originalURL string) *Account {
	return &Account{
		Name:        name,
		Email:       email,
		CreatedAt:   time.Now().UTC(),
		OriginalURL: originalURL,
	}
}
