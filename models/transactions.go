package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	Amount    float32   `json:"amount"`
	Operation string    `json:"operation"`
	Date      time.Time `json:"date"`
}

func (t Transaction) ToTransactionDTO() TransactionDTO {
	return TransactionDTO{
		ID:        t.ID,
		UserId:    t.UserId,
		Amount:    t.Amount,
		Operation: t.Operation,
		Date:      t.Date.Format(time.DateTime),
	}
}

type TransactionDTO struct {
	ID        int     `json:"id"`
	UserId    int     `json:"user_id" db:"user_id"`
	Amount    float32 `json:"amount"`
	Operation string  `json:"operation"`
	Date      string  `json:"date"`
}

func (t TransactionDTO) ToTransaction() (Transaction, error) {
	date, err := time.Parse(time.DateTime, t.Date)
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		ID:        t.ID,
		UserId:    t.UserId,
		Amount:    t.Amount,
		Operation: t.Operation,
		Date:      date,
	}, nil
}
