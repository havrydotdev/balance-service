package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gavrylenkoIvan/balance-service/internal/service"
	mock_service "github.com/gavrylenkoIvan/balance-service/internal/service/mocks"
	"github.com/gavrylenkoIvan/balance-service/models"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/gavrylenkoIvan/balance-service/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func returnFirstValue(values ...interface{}) interface{} {
	return values[0]
}

func parseTime(value string, t *testing.T) time.Time {
	timeAt, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Error(err)
	}

	return timeAt
}

func TestHandler_GetBalance(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, user int, currency string)

	testTable := []struct {
		name                 string
		userID               int
		currency             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:     "OK",
			userID:   1,
			currency: "EUR",
			mockBehavior: func(s *mock_service.MockUser, user int, currency string) {
				s.EXPECT().GetBalance(user, currency).Return(float32(4.13), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":4.13}`,
		},
		{
			name:     "OneMoreOK",
			userID:   2,
			currency: "EUR",
			mockBehavior: func(s *mock_service.MockUser, user int, currency string) {
				s.EXPECT().GetBalance(user, currency).Return(float32(32), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":2,"balance":32}`,
		},
		{
			name:     "OKinUAH",
			userID:   2,
			currency: "UAH",
			mockBehavior: func(s *mock_service.MockUser, user int, currency string) {
				uah, err := utils.Convert(32, "UAH")
				if err != nil {
					t.Error(err)
				}

				s.EXPECT().GetBalance(user, currency).Return(uah, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: fmt.Sprintf(`{"user_id":2,"balance":%.2f}`, returnFirstValue(utils.Convert(32, "UAH")).(float32)),
		},
		{
			name:     "NotValid",
			userID:   0,
			currency: "EUR",
			mockBehavior: func(s *mock_service.MockUser, user int, currency string) {
				s.EXPECT().GetBalance(user, currency).Return(float32(0), errors.New("user not found")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect user id"}`,
		},
		{
			name:     "DoesNotExist",
			userID:   400,
			currency: "EUR",
			mockBehavior: func(s *mock_service.MockUser, user int, currency string) {
				s.EXPECT().GetBalance(user, currency).Return(float32(0), errors.New("user not found"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"user not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.userID, testCase.currency)

			services := &service.Service{User: user}
			logger, err := logging.InitLogger()
			if err != nil {
				t.Error(err)
			}

			handler := NewHandler(services, logger)

			r := echo.New()
			r.GET("/balance/:user_id", handler.getBalance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET",
				fmt.Sprintf("/balance/%d?currency=%s", testCase.userID, testCase.currency), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_GetTransactions(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, userID int, page models.Page)

	testTable := []struct {
		name                 string
		userID               int
		page                 models.Page
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "OK",
			userID: 1,
			page: models.Page{
				Page:  1,
				Limit: 1,
				Sort:  "date",
			},
			mockBehavior: func(s *mock_service.MockUser, userID int, page models.Page) {
				s.EXPECT().GetTransactions(userID, page).Return([]models.Transaction{{
					ID:        1,
					UserId:    1,
					Amount:    30,
					Operation: "",
					Date:      parseTime("2023-06-10T00:13:35.315271Z", t),
				}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"user_id":1,"amount":30,"operation":"","date":"2023-06-10T00:13:35.315271Z"}]`,
		},
		{
			name:   "Multiple values + sort by ID",
			userID: 2,
			page: models.Page{
				Page:  1,
				Limit: 2,
				Sort:  "id",
			},
			mockBehavior: func(s *mock_service.MockUser, userID int, page models.Page) {
				s.EXPECT().GetTransactions(userID, page).Return([]models.Transaction{{
					ID:        2,
					UserId:    2,
					Amount:    101,
					Operation: "",
					Date:      parseTime("2023-06-10T00:13:35.315271Z", t),
				}, {
					ID:        3,
					UserId:    2,
					Amount:    32,
					Operation: "",
					Date:      parseTime("2023-06-10T00:13:35.315271Z", t),
				}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":2,"user_id":2,"amount":101,"operation":"","date":"2023-06-10T00:13:35.315271Z"},{"id":3,"user_id":2,"amount":32,"operation":"","date":"2023-06-10T00:13:35.315271Z"}]`,
		},
		{
			name:   "Multiple values + sort by ID + 2 page",
			userID: 3,
			page: models.Page{
				Page:  3,
				Limit: 2,
				Sort:  "id",
			},
			mockBehavior: func(s *mock_service.MockUser, userID int, page models.Page) {
				s.EXPECT().GetTransactions(userID, page).Return([]models.Transaction{{
					ID:        8,
					UserId:    3,
					Amount:    101,
					Operation: "",
					Date:      parseTime("2023-06-11T11:00:37.236094Z", t),
				}, {
					ID:        9,
					UserId:    3,
					Amount:    103,
					Operation: "",
					Date:      parseTime("2023-06-11T11:00:37.236094Z", t),
				}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":8,"user_id":3,"amount":101,"operation":"","date":"2023-06-11T11:00:37.236094Z"},{"id":9,"user_id":3,"amount":103,"operation":"","date":"2023-06-11T11:00:37.236094Z"}]`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.userID, testCase.page)

			services := &service.Service{User: user}
			logger, err := logging.InitLogger()
			if err != nil {
				t.Error(err)
			}

			handler := NewHandler(services, logger)

			r := echo.New()
			r.GET("/transactions/:user_id", handler.getTransactions)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET",
				fmt.Sprintf("/transactions/%d?page=%d&limit=%d&sort=%s",
					testCase.userID,
					testCase.page.Page,
					testCase.page.Limit,
					testCase.page.Sort), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_TopUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, input models.Input)

	testTable := []struct {
		name                 string
		input                models.Input
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: models.Input{
				UserId: 1,
				Amount: 30,
			},
			inputBody: `{"user_id":1,"amount":30}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().TopUp(input).Return(4.13+input.Amount, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":34.13}`,
		},
		{
			name: "Incorrect user id",
			input: models.Input{
				UserId: 0,
				Amount: 30,
			},
			inputBody: `{"user_id":0,"amount":30}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().TopUp(input).Return(float32(0), errors.New("user not found")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect user id"}`,
		},
		{
			name: "User does not exist",
			input: models.Input{
				UserId: 300,
				Amount: 30,
			},
			inputBody: `{"user_id":300,"amount":30}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().TopUp(input).Return(float32(0), errors.New("user not found")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"user not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			services := &service.Service{User: user}
			logger, err := logging.InitLogger()
			if err != nil {
				t.Error(err)
			}

			handler := NewHandler(services, logger)

			r := echo.New()
			r.POST("/top-up", handler.topUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST",
				"/top-up",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Add("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_Debit(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, input models.Input)

	testTable := []struct {
		name                 string
		input                models.Input
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: models.Input{
				UserId: 1,
				Amount: 1,
			},
			inputBody: `{"user_id":1,"amount":1}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().Debit(input).Return(4.13-input.Amount, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":3.13}`,
		},
		{
			name: "Incorrect user id",
			input: models.Input{
				UserId: 0,
				Amount: 30,
			},
			inputBody: `{"user_id":0,"amount":30}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().Debit(input).Return(float32(0), errors.New("user not found")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect user id"}`,
		},
		{
			name: "User does not exist",
			input: models.Input{
				UserId: 300,
				Amount: 30,
			},
			inputBody: `{"user_id":300,"amount":30}`,
			mockBehavior: func(s *mock_service.MockUser, input models.Input) {
				s.EXPECT().Debit(input).Return(float32(0), errors.New("user not found")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"user not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			services := &service.Service{User: user}
			logger, err := logging.InitLogger()
			if err != nil {
				t.Error(err)
			}

			handler := NewHandler(services, logger)

			r := echo.New()
			r.POST("/debit", handler.debit)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST",
				"/debit",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Add("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}

func TestHandler_Transfer(t *testing.T) {
	type mockBehavior func(s *mock_service.MockUser, input models.TransferInput)

	testTable := []struct {
		name                 string
		input                models.TransferInput
		inputBody            string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "OK",
			input: models.TransferInput{
				UserId: 1,
				ToId:   2,
				Amount: 4.13,
			},
			inputBody: `{"user_id":1,"to_id":2,"amount":4.13}`,
			mockBehavior: func(s *mock_service.MockUser, input models.TransferInput) {
				s.EXPECT().Transfer(input).Return(4.13-input.Amount, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":0}`,
		},
		{
			name: "Incorrect user id",
			input: models.TransferInput{
				UserId: 0,
				ToId:   2,
				Amount: 4.13,
			},
			inputBody: `{"user_id":0,"to_id":2,"amount":4.13}`,
			mockBehavior: func(s *mock_service.MockUser, input models.TransferInput) {
				s.EXPECT().Transfer(input).Return(4.13-input.Amount, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect user id"}`,
		},
		{
			name: "Incorrect to id",
			input: models.TransferInput{
				UserId: 1,
				ToId:   0,
				Amount: 4.13,
			},
			inputBody: `{"user_id":1,"to_id":0,"amount":4.13}`,
			mockBehavior: func(s *mock_service.MockUser, input models.TransferInput) {
				s.EXPECT().Transfer(input).Return(4.13-input.Amount, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect to id"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			user := mock_service.NewMockUser(c)
			testCase.mockBehavior(user, testCase.input)

			services := &service.Service{User: user}
			logger, err := logging.InitLogger()
			if err != nil {
				t.Error(err)
			}

			handler := NewHandler(services, logger)

			r := echo.New()
			r.POST("/transfer", handler.transfer)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST",
				"/transfer",
				bytes.NewBufferString(testCase.inputBody))
			req.Header.Add("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}
