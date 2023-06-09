package service

import (
	"github.com/gavrylenkoIvan/balance-service/internal/repo"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
)

type Service struct {
	User
}

func NewService(repo *repo.Repo, log logging.Logger) *Service {
	return &Service{
		User: NewUserService(repo, log),
	}
}
