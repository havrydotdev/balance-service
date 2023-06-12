package repo

//
//import (
//	"errors"
//	"fmt"
//	"github.com/gavrylenkoIvan/balance-service/models"
//	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
//	"github.com/stretchr/testify/assert"
//	sqlmock "github.com/zhashkevych/go-sqlxmock"
//	"testing"
//)
//
//func TestUserRepository_Balance(t *testing.T) {
//	db, mock, err := sqlmock.Newx()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	logger, err := logging.InitLogger()
//	if err != nil {
//		t.Error(err)
//	}
//
//	r := NewUserRepo(db, logger)
//
//	type mockBehavior func(userId int)
//
//	tests := []struct {
//		name    string
//		mock    mockBehavior
//		input   int
//		want    *models.User
//		wantErr bool
//	}{
//		{
//			name: "Ok",
//			mock: func(userId int) {
//				rows := sqlmock.NewRows([]string{"id", "balance"}).AddRow(userId, 4.13)
//				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
//					WithArgs(userId).WillReturnRows(rows)
//			},
//			input: 1,
//			want: &models.User{
//				ID:      1,
//				Balance: 10,
//			},
//		},
//		{
//			name: "User has no balance",
//			mock: func(userId int) {
//				rows := sqlmock.NewRows([]string{"id", "user_id", "balance"}).AddRow(0, 0, 0).RowError(0, errors.New("some error"))
//				mock.ExpectQuery(fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)).
//					WithArgs(userId).WillReturnRows(rows)
//			},
//			input:   0,
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.mock(tt.input)
//
//			got, err := r.GetBalance(tt.input)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//				assert.Equal(t, tt.want, got)
//			}
//			assert.NoError(t, mock.ExpectationsWereMet())
//		})
//	}
//}
