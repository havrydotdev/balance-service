package repo

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	usersTable        = "users"
	transactionsTable = "transactions"
)

type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
	SSL      string
}

func InitDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Name, cfg.Password, cfg.SSL))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(0)

	return db, nil
}
