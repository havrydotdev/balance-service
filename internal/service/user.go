package service

import (
	"github.com/gavrylenkoIvan/balance-service/internal/repo"
	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/gavrylenkoIvan/balance-service/pkg/utils"
)

type UserService struct {
	repo repo.User
	log  logging.Logger
}

func NewUserService(repo repo.User, log logging.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

func (s *UserService) TopUp(input models.Input) (float32, error) {
	return s.repo.TopUp(input)
}

func (s *UserService) Debit(input models.Input) (float32, error) {
	return s.repo.Debit(input)
}

func (s *UserService) Transfer(input models.TransferInput) (float32, error) {
	return s.repo.Transfer(input)
}

func (s *UserService) GetTransactions(id int, page models.Page) ([]models.Transaction, error) {
	return s.repo.GetTransactions(id, page)
}

func (s *UserService) GetBalance(id int, currency string) (float32, error) {
	euro, err := s.repo.GetBalance(id)
	if err != nil {
		return 0, err
	}

	return utils.Convert(euro, currency)
}
