package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/gavrylenkoIvan/balance-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserRepository_GetBalance(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger, err := logging.InitLogger()
	if err != nil {
		t.Error(err)
	}

	r := NewUserRepo(sqlxDB, logger)

	type mockBehavior func(userID int)

	tests := []struct {
		name      string
		mock      mockBehavior
		userID    int
		want      float32
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Ok",
			mock: func(userID int) {
				rows := sqlmock.NewRows([]string{"balance"}).AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(userID).WillReturnRows(rows)
			},
			userID:  1,
			want:    10,
			wantErr: false,
		},
		{
			name: "Error no rows",
			mock: func(userID int) {
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(userID).WillReturnError(sql.ErrNoRows)
			},
			userID:    100,
			want:      0,
			wantErr:   true,
			wantedErr: "user not found",
		},
		{
			name: "Random error",
			mock: func(userID int) {
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(userID).WillReturnError(errors.New("db is not valid"))
			},
			userID:    100,
			want:      0,
			wantErr:   true,
			wantedErr: "db is not valid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.userID)

			got, err := r.GetBalance(tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, errors.New(tt.wantedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetTransactions(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger, err := logging.InitLogger()
	if err != nil {
		t.Error(err)
	}

	r := NewUserRepo(sqlxDB, logger)

	type mockBehavior func(userID int)

	tests := []struct {
		name      string
		mock      mockBehavior
		userID    int
		page      models.Page
		want      []models.Transaction
		wantErr   bool
		wantedErr error
	}{
		{
			name: "Ok",
			mock: func(userID int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).
					AddRow(1, 1, 10, "Debit by transfer 10EUR", time.Now().Format(time.DateTime)).
					AddRow(2, 1, 5, "Top-up by transfer 5EUR", time.Now().Format(time.DateTime))
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)", transactionsTable)).
					WithArgs(userID).WillReturnRows(rows)
			},
			userID: 1,
			want: []models.Transaction{
				{
					ID:        1,
					UserId:    1,
					Amount:    10,
					Operation: "Debit by transfer 10EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
				{
					ID:        2,
					UserId:    1,
					Amount:    5,
					Operation: "Top-up by transfer 5EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
			},
			page: models.Page{
				Page:  1,
				Limit: 10,
				Sort:  "date",
			},
			wantErr: false,
		},
		{
			name: "User does not exist",
			mock: func(userID int) {
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)", transactionsTable)).
					WithArgs(userID).WillReturnError(sql.ErrNoRows)
			},
			want: []models.Transaction{
				{
					ID:        1,
					UserId:    1,
					Amount:    10,
					Operation: "Debit by transfer 10EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
				{
					ID:        2,
					UserId:    1,
					Amount:    5,
					Operation: "Top-up by transfer 5EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
			},
			page: models.Page{
				Page:  1,
				Limit: 10,
				Sort:  "date",
			},
			wantErr:   true,
			wantedErr: errors.New("user not found"),
		},
		{
			name: "Random error",
			mock: func(userID int) {
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)", transactionsTable)).
					WithArgs(userID).WillReturnError(errors.New("db is not valid"))
			},
			want: []models.Transaction{
				{
					ID:        1,
					UserId:    1,
					Amount:    10,
					Operation: "Debit by transfer 10EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
				{
					ID:        2,
					UserId:    1,
					Amount:    5,
					Operation: "Top-up by transfer 5EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
			},
			page: models.Page{
				Page:  1,
				Limit: 10,
				Sort:  "date",
			},
			wantErr:   true,
			wantedErr: errors.New("db is not valid"),
		},
		{
			name: "Failed to convert date",
			mock: func(userID int) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "amount", "operation", "date"}).
					AddRow(1, 2, 100, "", "1849q9")

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+) ORDER BY (.+) LIMIT (.+) OFFSET (.+)", transactionsTable)).
					WithArgs(userID).WillReturnRows(rows)
			},
			want: []models.Transaction{
				{
					ID:        1,
					UserId:    1,
					Amount:    10,
					Operation: "Debit by transfer 10EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
				{
					ID:        2,
					UserId:    1,
					Amount:    5,
					Operation: "Top-up by transfer 5EUR",
					Date:      utils.ParseTime(time.Now().Format(time.DateTime), t),
				},
			},
			page: models.Page{
				Page:  1,
				Limit: 10,
				Sort:  "date",
			},
			wantErr: true,
			wantedErr: &time.ParseError{
				Layout:     "2006-01-02 15:04:05",
				Value:      "1849q9",
				LayoutElem: "-",
				ValueElem:  "q9",
				Message:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.userID)

			got, err := r.GetTransactions(tt.userID, tt.page)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, tt.wantedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_TopUp(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger, err := logging.InitLogger()
	if err != nil {
		t.Error(err)
	}

	r := NewUserRepo(sqlxDB, logger)

	type mockBehavior func(input models.Input)

	tests := []struct {
		name      string
		mock      mockBehavior
		want      float32
		input     models.Input
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Ok",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date := time.Now().Format("01-02-2006 15:04:05")
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Top-up by bank_card %fEUR", input.Amount), date).
					WillReturnResult(result)

				mock.ExpectCommit()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "Failed to insert",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date := time.Now().Format("01-02-2006 15:04:05")
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Top-up by bank_card %fEUR", input.Amount), date).
					WillReturnResult(result)

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to insert new transaction, rollback",
		},
		{
			name: "User does not exist",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 0))

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "user not found",
		},
		{
			name: "Failed to begin tx",
			mock: func(input models.Input) {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin tx"))
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to begin tx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.TopUp(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, errors.New(tt.wantedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Debit(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger, err := logging.InitLogger()
	if err != nil {
		t.Error(err)
	}

	r := NewUserRepo(sqlxDB, logger)

	type mockBehavior func(input models.Input)

	tests := []struct {
		name      string
		mock      mockBehavior
		want      float32
		input     models.Input
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Ok",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date := time.Now().Format("01-02-2006 15:04:05")
				result := sqlmock.NewResult(1, 1)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by purchase %fEUR", input.Amount), date).
					WillReturnResult(result)

				mock.ExpectCommit()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "Failed to insert",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date := time.Now().Format("01-02-2006 15:04:05")
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by purchase %fEUR", input.Amount), date).
					WillReturnResult(result)

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to insert new transaction, rollback",
		},
		{
			name: "User does not exist",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 0))

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "user not found",
		},
		{
			name: "Failed to begin tx",
			mock: func(input models.Input) {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin tx"))
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to begin tx",
		},
		{
			name: "Select returned error",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnError(errors.New("failed to connect to db"))

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to connect to db",
		},
		{
			name: "No enough money",
			mock: func(input models.Input) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectRollback()
			},
			input: models.Input{
				UserId: 1,
				Amount: 11,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "not enough money to perform purchase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Debit(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, errors.New(tt.wantedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Transfer(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	defer mockDB.Close()
	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	logger, err := logging.InitLogger()
	if err != nil {
		t.Error(err)
	}

	r := NewUserRepo(sqlxDB, logger)

	type mockBehavior func(input models.TransferInput)

	tests := []struct {
		name      string
		mock      mockBehavior
		want      float32
		input     models.TransferInput
		wantErr   bool
		wantedErr string
	}{
		{
			name: "Ok",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin()

				selectRows2 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows2)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date2 := time.Now().Format("01-02-2006 15:04:05")
				result2 := sqlmock.NewResult(1, 1)

				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by transfer %fEUR", input.Amount), date2).
					WillReturnResult(result2)

				selectRows1 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.ToId).
					WillReturnRows(selectRows1)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.ToId).WillReturnResult(sqlmock.NewResult(1, 1))

				date1 := time.Now().Format("01-02-2006 15:04:05")
				result1 := sqlmock.NewResult(1, 1)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.ToId, input.Amount, fmt.Sprintf("Top-up by transfer %fEUR", input.Amount), date1).
					WillReturnResult(result1)

				mock.ExpectCommit()
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:    20,
			wantErr: false,
		},
		{
			name: "Failed to insert debit transaction",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date := time.Now().Format("01-02-2006 15:04:05")
				result := sqlmock.NewResult(0, 0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by transfer %fEUR", input.Amount), date).
					WillReturnResult(result)

				mock.ExpectRollback()
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to insert new transaction, rollback",
		},
		{
			name: "Failed to insert top-up transaction",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin()

				selectRows2 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows2)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date2 := time.Now().Format("01-02-2006 15:04:05")
				result2 := sqlmock.NewResult(1, 1)

				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by transfer %fEUR", input.Amount), date2).
					WillReturnResult(result2)

				selectRows1 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.ToId).
					WillReturnRows(selectRows1)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.ToId).WillReturnResult(sqlmock.NewResult(1, 1))

				date1 := time.Now().Format("01-02-2006 15:04:05")
				result1 := sqlmock.NewResult(1, 0)
				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.ToId, input.Amount, fmt.Sprintf("Top-up by transfer %fEUR", input.Amount), date1).
					WillReturnResult(result1)

				mock.ExpectRollback()
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to insert new transaction, rollback",
		},
		{
			name: "User does not exist",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin()

				selectRows := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)
				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 0))

				mock.ExpectRollback()
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "user not found",
		},
		{
			name: "Receiver does not exist",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin()

				selectRows2 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.UserId).
					WillReturnRows(selectRows2)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.UserId).WillReturnResult(sqlmock.NewResult(1, 1))

				date2 := time.Now().Format("01-02-2006 15:04:05")
				result2 := sqlmock.NewResult(1, 1)

				mock.ExpectExec(fmt.Sprintf("INSERT INTO %s", transactionsTable)).
					WithArgs(input.UserId, input.Amount, fmt.Sprintf("Debit by transfer %fEUR", input.Amount), date2).
					WillReturnResult(result2)

				selectRows1 := sqlmock.NewRows([]string{"balance"}).
					AddRow(10)

				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
					WithArgs(input.ToId).
					WillReturnRows(selectRows1)

				mock.ExpectExec(fmt.Sprintf("UPDATE %s SET (.+) WHERE (.+) RETURNING (.+)", usersTable)).
					WithArgs(input.ToId).WillReturnResult(sqlmock.NewResult(1, 0))

				mock.ExpectRollback()
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "user not found",
		},
		{
			name: "Failed to begin tx",
			mock: func(input models.TransferInput) {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin tx"))
			},
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 10,
			},
			want:      0,
			wantErr:   true,
			wantedErr: "failed to begin tx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input)

			got, err := r.Transfer(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, err, errors.New(tt.wantedErr))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
