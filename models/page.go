package models

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type Page struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func PageFromRequest(c echo.Context) Page {
	page := 1
	limit := 10
	sort := "date"
	if p := c.QueryParam("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}

	if l := c.QueryParam("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	if s := c.QueryParam("sort"); s != "" {
		sort = s
	}

	return Page{
		Page:  page,
		Limit: limit,
		Sort:  sort,
	}
}
