package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/jmoiron/sqlx"
	"time"
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
	}, fmt.Sprintf("Top-up by transfer %fEUR", input.Amount), tx)
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
	return r.ChangeBalance(input, "-", operation, tx)
}

func (r *UserRepo) ChangeBalance(input models.Input, action string, operation string, tx *sql.Tx) (float32, error) {
	var balance float32

	check := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&balance, check, input.UserId)
	if err != nil {
		return 0, err
	}

	if balance-input.Amount < 0 && action == "-" {
		return 0, errors.New("not enough money to perform purchase")
	}

	query := fmt.Sprintf("UPDATE %s SET balance = balance %s %f WHERE id = $1 RETURNING balance",
		usersTable, action, input.Amount)

	res, err := tx.Exec(query, input.UserId)
	if err != nil {
		return 0, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	if affected == 0 {
		return 0, errors.New("user not found")
	}

	insert := fmt.Sprintf("INSERT INTO %s (user_id, amount, operation, date) VALUES ($1, $2, $3, $4)",
		transactionsTable)

	result, err := tx.Exec(insert, input.UserId, input.Amount, operation, time.Now().Format("01-02-2006 15:04:05"))
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if count == 0 {
		return 0, errors.New("failed to insert new transaction, rollback")
	}

	if action == "+" {
		return balance + input.Amount, nil
	} else {
		return balance - input.Amount, nil
	}
}

func (r *UserRepo) GetTransactions(id int, page models.Page) ([]models.Transaction, error) {
	var transactions []models.TransactionDTO
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY %s LIMIT %d OFFSET %d",
		transactionsTable, page.Sort, page.Limit, (page.Page-1)*page.Limit)

	err := r.db.Select(&transactions, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	result := make([]models.Transaction, 0, len(transactions))
	for i := 0; i < len(transactions); i++ {
		t, err := time.Parse(time.DateTime, transactions[i].Date)
		if err != nil {
			return nil, err
		}

		result = append(result, models.Transaction{
			ID:        transactions[i].ID,
			UserId:    transactions[i].UserId,
			Amount:    transactions[i].Amount,
			Operation: transactions[i].Operation,
			Date:      t,
		})
	}

	return result, nil
}

func (r *UserRepo) GetBalance(id int) (float32, error) {
	var balance float32
	query := fmt.Sprintf("SELECT balance FROM %s WHERE id = $1", usersTable)
	err := r.db.Get(&balance, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("user not found")
		}

		return 0, err
	}

	r.log.LogRepo("GET", "GetBalance", true, balance)
	return balance, nil
}
