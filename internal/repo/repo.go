package repo

import (
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	User
}

func NewRepo(db *sqlx.DB, log logging.Logger) *Repo {
	return &Repo{
		User: NewUserRepo(db, log),
	}
}
