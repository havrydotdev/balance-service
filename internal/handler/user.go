package handler

import (
	"net/http"
	"strconv"

	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/labstack/echo/v4"
)

func (h *Handler) transfer(c echo.Context) error {
	var input models.TransferInput
	if err := c.Bind(&input); err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusBadRequest, c)
		return err
	}

	balance, err := h.s.Transfer(input)
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusInternalServerError, c)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
		"user_id": input.UserId,
	})
}

func (h *Handler) debit(c echo.Context) error {
	var input models.Input
	if err := c.Bind(&input); err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusBadRequest, c)
		return err
	}

	balance, err := h.s.Debit(input)
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusInternalServerError, c)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
		"user_id": input.UserId,
	})
}

func (h *Handler) topUp(c echo.Context) error {
	var input models.Input
	if err := c.Bind(&input); err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusBadRequest, c)
		return err
	}

	balance, err := h.s.TopUp(input)
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusInternalServerError, c)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"balance": balance,
		"user_id": input.UserId,
	})
}

func (h *Handler) getBalance(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusBadRequest, c)
		return err
	}

	balance, err := h.s.GetBalance(userId, c.QueryParam("currency"))
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusInternalServerError, c)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": userId,
		"balance": balance,
	})
}

func (h *Handler) getTransactions(c echo.Context) error {
	userId, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusBadRequest, c)
		return err
	}

	page := models.PageFromRequest(c)

	transactions, err := h.s.GetTransactions(userId, page)
	if err != nil {
		h.log.ErrorResponse(err.Error(), http.StatusInternalServerError, c)
		return err
	}

	return c.JSON(http.StatusOK, transactions)
}
