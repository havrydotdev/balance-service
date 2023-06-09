package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/jmoiron/sqlx"
)

type User interface {
	GetBalance(id int) (float32, error)
	GetTransactions(id int, page models.Page) ([]models.Transaction, error)
	DecreaseBalanceTx(input models.Input, operation string, tx *sql.Tx) (float32, error)
	IncreaseBalanceTx(input models.Input, operation string, tx *sql.Tx) (float32, error)
	TopUp(input models.Input) (float32, error)
	Debit(input models.Input) (float32, error)
	Transfer(input models.TransferInput) (float32, error)
	ChangeBalance(input models.Input, action string, operation string, tx *sql.Tx) (float32, error)
}

type UserRepo struct {
	db  *sqlx.DB
	log logging.Logger
}

func NewUserRepo(db *sqlx.DB, log logging.Logger) *UserRepo {
	return &UserRepo{
		db:  db,
		log: log,
	}
}

func (r *UserRepo) Transfer(input models.TransferInput) (float32, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	_, err = r.DecreaseBalanceTx(models.Input{
		UserId: input.UserId,
		Amount: input.Amount,
	}, fmt.Sprintf("Debit by transfer %fEUR", input.Amount), tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	balance, err := r.IncreaseBalanceTx(models.Input{
		UserId: input.ToId,
		Amount: input.Amount,
	}, fmt.Sprintf("Top-up by transfer %f EUR", input.Amount), tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return balance, tx.Commit()
}

func (r *UserRepo) Debit(input models.Input) (float32, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	res, err := r.DecreaseBalanceTx(input, fmt.Sprintf("Debit by purchase %fEUR", input.Amount), tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return res, tx.Commit()
}

func (r *UserRepo) TopUp(input models.Input) (float32, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	balance, err := r.IncreaseBalanceTx(input, fmt.Sprintf("Top-up by bank_card %fEUR", input.Amount), tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return balance, tx.Commit()
}

func (r *UserRepo) IncreaseBalanceTx(input models.Input, operation string, tx *sql.Tx) (float32, error) {
	return r.ChangeBalance(input, "+", operation, tx)
}

func (r *UserRepo) DecreaseBalanceTx(input models.Input, operation string, tx *sql.Tx) (float32, error) {
	var balance float32
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&balance, query, input.UserId)
	if err != nil {
		return 0, err
	}

	if balance-input.Amount < 0 {
		return 0, errors.New("not enough money to perform purchase")
	}

	return r.ChangeBalance(input, "-", operation, tx)
}

func (r *UserRepo) ChangeBalance(input models.Input, action string, operation string, tx *sql.Tx) (float32, error) {
	query := fmt.Sprintf("UPDATE %s SET balance = balance %s %f WHERE id = $1 RETURNING balance",
		usersTable, action, input.Amount)

	res := tx.QueryRow(query, input.UserId)
	if res.Err() != nil {
		tx.Rollback()
		return 0, res.Err()
	}

	var balance float32
	err := res.Scan(&balance)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	insert := fmt.Sprintf("INSERT INTO %s (user_id, amount, operation, date) VALUES ($1, $2, $3, $4)",
		transactionsTable)

	_, err = tx.Exec(insert, input.UserId, input.Amount, operation, time.Now())
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return balance, nil
}

func (r *UserRepo) GetTransactions(id int, page models.Page) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY %s LIMIT %d OFFSET %d",
		transactionsTable, page.Sort, page.Limit, (page.Page-1)*page.Limit)

	err := r.db.Select(&transactions, query, id)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *UserRepo) GetBalance(id int) (float32, error) {
	var balance float32
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&balance, query, id)
	if err != nil {
		return 0, err
	}

	r.log.LogRepo("GET", "GetBalance", true, balance)
	return balance, nil
}
