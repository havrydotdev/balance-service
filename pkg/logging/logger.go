package logging

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Logger interface {
	Fatal(msg string)
	Fatalf(msgf string, args ...interface{})
	Info(msg string)
	Infof(msgf string, args ...interface{})
	ErrorResponse(msg string, err int, c echo.Context) error
	Infow(fields Fields)
	LogRepo(method, info string, ok bool, resp interface{})
}

type logger struct {
	logger *zap.SugaredLogger
}

type ErrorResponse struct {
	Message string `json:"msg"`
}

type Fields map[string]interface{}

func ln(msg string) string {
	return fmt.Sprintf("%s\n", msg)
}

func InitLogger() (*logger, error) {
	l, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	return &logger{
		logger: l.Sugar(),
	}, nil
}

func (l *logger) Fatal(msg string) {
	l.logger.Fatal(ln(msg))
}

func (l *logger) Fatalf(msgf string, args ...interface{}) {
	l.logger.Fatalf(ln(msgf), args...)
}

func (l *logger) Info(msg string) {
	l.logger.Info(ln(msg))
}

func (l *logger) Infof(msgf string, args ...interface{}) {
	l.logger.Infof(ln(msgf), args...)
}

func (l *logger) ErrorResponse(msg string, err int, c echo.Context) error {
	l.logger.Error(msg)
	return c.JSON(err, ErrorResponse{
		Message: msg,
	})
}

func (l *logger) Infow(fields Fields) {
	res, err := json.Marshal(fields)
	if err != nil {
		l.logger.Error(err.Error())
		return
	}

	l.logger.Info(string(res))
}

func (l *logger) LogRepo(method, info string, ok bool, resp interface{}) {
	l.Infow(Fields{
		"method": method,
		"func":   info,
		"ok":     ok,
		"resp":   resp,
	})
}
