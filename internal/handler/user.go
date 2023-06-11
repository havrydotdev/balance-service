package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/labstack/echo/v4"
)

type transactionResponse struct {
	UserId  int     `json:"user_id"`
	Balance float32 `json:"balance"`
}

// @Summary Transfer money
// @Tags balance
// @Description Transfer money from one user to another
// @ID transfer
// @Accept  json
// @Produce  json
// @Param input body models.TransferInput true "transfer info"
// @Success 200 {object} transactionResponse
// @Failure 400,404 {object} logging.ErrorResponse
// @Failure 500 {object} logging.ErrorResponse
// @Failure default {object} logging.ErrorResponse
// @Router /transfer [post]
func (h *Handler) transfer(c echo.Context) error {
	var input models.TransferInput
	if err := c.Bind(&input); err != nil {
		return h.log.ErrorResponse(http.StatusBadRequest, err)
	}

	if input.UserId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect user id"))
	} else if input.ToId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect to id"))
	}

	balance, err := h.s.Transfer(input)
	if err != nil {
		return h.log.ErrorResponse(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, transactionResponse{
		UserId:  input.UserId,
		Balance: balance,
	})
}

// @Summary Debit from card
// @Tags balance
// @Description Decreases user`s balance by input.Amount
// @ID debit
// @Accept  json
// @Produce  json
// @Param input body models.Input true "debit input"
// @Success 200 {object} transactionResponse
// @Failure 400,404 {object} logging.ErrorResponse
// @Failure 500 {object} logging.ErrorResponse
// @Failure default {object} logging.ErrorResponse
// @Router /debit [post]
func (h *Handler) debit(c echo.Context) error {
	var input models.Input
	if err := c.Bind(&input); err != nil {
		return h.log.ErrorResponse(http.StatusBadRequest, err)
	}

	if input.UserId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect user id"))
	}

	balance, err := h.s.Debit(input)
	if err != nil {
		return h.log.ErrorResponse(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, transactionResponse{
		UserId:  input.UserId,
		Balance: balance,
	})
}

// @Summary Top up
// @Tags balance
// @Description Increases user`s balance by input.Amount
// @ID top-up
// @Accept  json
// @Produce  json
// @Param input body models.Input true "top up input"
// @Success 200 {object} transactionResponse
// @Failure 400,404 {object} logging.ErrorResponse
// @Failure 500 {object} logging.ErrorResponse
// @Failure default {object} logging.ErrorResponse
// @Router /top-up [post]
func (h *Handler) topUp(c echo.Context) error {
	var input models.Input
	if err := c.Bind(&input); err != nil {
		return h.log.ErrorResponse(http.StatusBadRequest, err)
	}

	if input.UserId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect user id"))
	}

	balance, err := h.s.TopUp(input)
	if err != nil {
		return h.log.ErrorResponse(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, transactionResponse{
		UserId:  input.UserId,
		Balance: balance,
	})
}

// @Summary Get balance
// @Tags balance
// @Description Returns user`s balance
// @ID get-balance
// @Produce  json
// @Param        id   path      int  true  "User ID"
// @Success 200 {object} transactionResponse
// @Failure 400,404 {object} logging.ErrorResponse
// @Failure 500 {object} logging.ErrorResponse
// @Failure default {object} logging.ErrorResponse
// @Router /balance/{id} [get]
func (h *Handler) getBalance(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return h.log.ErrorResponse(http.StatusBadRequest, err)
	}

	if userId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect user id"))
	}

	balance, err := h.s.GetBalance(userId, c.QueryParam("currency"))
	if err != nil {
		return h.log.ErrorResponse(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, transactionResponse{
		UserId:  userId,
		Balance: balance,
	})
}

// @Summary Get transactions
// @Tags balance
// @Description Returns user`s transactions
// @ID get-transactions
// @Produce  json
// @Param        id   path      int  true  "User ID"
// @Success 200 {object} []models.Transaction
// @Failure 400,404 {object} logging.ErrorResponse
// @Failure 500 {object} logging.ErrorResponse
// @Failure default {object} logging.ErrorResponse
// @Router /transactions/{id} [get]
func (h *Handler) getTransactions(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		return h.log.ErrorResponse(http.StatusBadRequest, err)
	}

	if userId <= 0 {
		return h.log.ErrorResponse(http.StatusBadRequest, errors.New("incorrect user id"))
	}

	page := models.PageFromRequest(c)

	transactions, err := h.s.GetTransactions(userId, page)
	if err != nil {
		return h.log.ErrorResponse(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, transactions)
}
