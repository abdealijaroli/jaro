package store

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/abdealijaroli/jaro/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	CreateWaitlist(string, string) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	connStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	for _, initFunc := range []func() error{
		s.CreateAccountTable,
		s.CreateShortURLTable,
		s.CreateWaitlistTable,
	} {
		if err := initFunc(); err != nil {
			return fmt.Errorf("initialization error: %w", err)
		}
	}
	return nil
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}

func (s *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
        id SERIAL PRIMARY KEY,
        name VARCHAR(50),
        email VARCHAR(50),
		password_hash VARCHAR(255),
        created_at TIMESTAMP,
        short_url VARCHAR(255),
        original_url VARCHAR(255)
    )`

	_, err := s.db.Exec(query)

	return err
}

func (s *PostgresStore) CreateShortURLTable() error {
	query := `CREATE TABLE IF NOT EXISTS short_urls (
		id SERIAL PRIMARY KEY,
		original_url VARCHAR(255),
		short_url VARCHAR(10) UNIQUE,	
		created_at TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateWaitlistTable() error {
	query := `CREATE TABLE IF NOT EXISTS waitlist (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50),
		email VARCHAR(50),
		created_at TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateWaitlist(name, email string) error {
	insertQuery := `INSERT INTO waitlist (name, email, created_at) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(insertQuery, name, email, time.Now())
	return err
}

func (s *PostgresStore) AddShortURLToDB(originalURL, shortURL string) error {
	query := `INSERT INTO short_urls (original_url, short_url, created_at) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(query, originalURL, shortURL, time.Now())
	return err
}

func (s *PostgresStore) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM short_urls WHERE short_url = $1`
	err := s.db.QueryRow(query, shortURL).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("short URL not found")
		}
		return "", err
	}
	return originalURL, nil
}

func (s *PostgresStore) CreateAccount(acc *types.Account) error {
	query := `INSERT INTO accounts (name, email, created_at, short_url, original_url) VALUES ($1, $2, $3, $4, $5)`
	err := s.db.QueryRow(query, acc.Name, acc.Email, acc.CreatedAt, acc.ShortURL, acc.OriginalURL).Scan(&acc.ID)
	return err
}

func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) UpdateAccount(acc *types.Account) error {
	query := `UPDATE accounts SET name = $1, email = $2, short_url = $3, original_url = $4 WHERE id = $5`
	_, err := s.db.Exec(query, acc.Name, acc.Email, acc.ShortURL, acc.OriginalURL, acc.ID)
	return err
}

func (s *PostgresStore) GetAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*types.Account{}

	for rows.Next() {
		acc, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (s *PostgresStore) GetAccountByID(id int) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func scanIntoAccount(rows *sql.Rows) (*types.Account, error) {
	acc := &types.Account{}
	err := rows.Scan(&acc.ID, &acc.Name, &acc.Email, &acc.CreatedAt, &acc.ShortURL, &acc.OriginalURL)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
