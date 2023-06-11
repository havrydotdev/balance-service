package handler

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gavrylenkoIvan/balance-service/internal/service"
	mock_service "github.com/gavrylenkoIvan/balance-service/internal/service/mocks"
	"github.com/gavrylenkoIvan/balance-service/pkg/logging"
	"github.com/gavrylenkoIvan/balance-service/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func returnFirstValue(values ...interface{}) interface{} {
	return values[0]
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
				s.EXPECT().GetBalance(user, currency).Return(float32(0), errors.New("user not found"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"user not found"}`,
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
	type mockBehavior func(s *mock_service.MockUser, user int)

	testTable := []struct {
		name                 string
		userID               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "OK",
			userID: 1,
			mockBehavior: func(s *mock_service.MockUser, user int) {
				s.EXPECT().GetBalance(1, "").Return(float32(4.13), nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"user_id":1,"balance":4.13}`,
		},
		{
			name:   "NotValid",
			userID: 0,
			mockBehavior: func(s *mock_service.MockUser, user int) {
				s.EXPECT().GetBalance(0, "").Return(float32(0), errors.New("user not found"))
			},
			expectedStatusCode:   404,
			expectedResponseBody: `{"message":"user not found"}`,
		},
		{
			name:   "DoesNotExist",
			userID: 400,
			mockBehavior: func(s *mock_service.MockUser, user int) {
				s.EXPECT().GetBalance(400, "").Return(float32(0), errors.New("user not found"))
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
			testCase.mockBehavior(user, testCase.userID)

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
				fmt.Sprintf("/balance/%d", testCase.userID), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, strings.ReplaceAll(w.Body.String(), "\n", ""))
		})
	}
}
