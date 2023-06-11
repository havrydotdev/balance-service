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
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetBalance(t *testing.T) {
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
