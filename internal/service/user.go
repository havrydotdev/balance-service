package service

import (
	"encoding/json"
	"net/http"

	"github.com/gavrylenkoIvan/balance-service/internal/repo"
	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
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

	return convert(euro, currency)
}

func convert(euro float32, currency string) (float32, error) {
	if currency == "" {
		return euro, nil
	}

	resp, err := http.Get("http://api.exchangeratesapi.io/v1/latest?access_key=5bb179314fdbfaa6a839358e571d426f&base=EUR&symbols=" + currency)
	if err != nil {
		return 0, err
	}

	var get models.Response
	json.NewDecoder(resp.Body).Decode(&get)
	return euro * get.Rates[currency], nil
}
