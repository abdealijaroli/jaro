package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	connStr := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO accounts (name, email, created_at, short_url, original_url) VALUES ($1, $2, $3, $4, $5)`
	err := s.db.QueryRow(query, acc.Name, acc.Email, acc.CreatedAt, acc.ShortURL, acc.OriginalURL).Scan(&acc.ID)
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) UpdateAccount(acc *Account) error {
	query := `UPDATE accounts SET name = $1, email = $2, short_url = $3, original_url = $4 WHERE id = $5`
	_, err := s.db.Exec(query, acc.Name, acc.Email, acc.ShortURL, acc.OriginalURL, acc.ID)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*Account{}

	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	acc := &Account{}
	err := rows.Scan(&acc.ID, &acc.Name, &acc.Email, &acc.CreatedAt, &acc.ShortURL, &acc.OriginalURL)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
