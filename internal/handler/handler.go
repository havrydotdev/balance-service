package handler

import (
	"github.com/gavrylenkoIvan/balance-service/internal/service"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	s   *service.Service
	log logging.Logger
}

func NewHandler(s *service.Service, log logging.Logger) *Handler {
	return &Handler{
		s:   s,
		log: log,
	}
}

func (h *Handler) InitRoutes() *echo.Echo {
	r := echo.New()

	r.Use(middleware.Logger())
	r.GET("/balance/:user_id", h.getBalance)
	r.GET("/transactions/:user_id", h.getTransactions)
	r.POST("/top-up", h.topUp)
	r.POST("/debit", h.debit)
	r.POST("/transfer", h.transfer)
	return r
}
