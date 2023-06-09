package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Amount    float32   `json:"amount"`
	Operation string    `json:"operation"`
	Date      time.Time `json:"date"`
}

type TransactionDTO struct {
	UserId    int       `json:"user_id" db:"user_id"`
	Amount    float32   `json:"amount"`
	Operation string    `json:"operation"`
	Date      time.Time `json:"date"`
}
