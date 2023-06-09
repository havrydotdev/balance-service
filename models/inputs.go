package models

type Input struct {
	UserId int     `json:"user_id"`
	Amount float32 `json:"amount"`
}

type TransferInput struct {
	ToId   int     `json:"to_id"`
	UserId int     `json:"user_id"`
	Amount float32 `json:"amount"`
}
