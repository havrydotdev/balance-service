package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gavrylenkoIvan/balance-service/models"
)

func Convert(euro float32, currency string) (float32, error) {
	if currency == "" {
		return euro, nil
	}

	resp, err := http.Get("http://api.exchangeratesapi.io/v1/latest?access_key=5bb179314fdbfaa6a839358e571d426f&base=EUR&symbols=" + currency)
	if err != nil {
		return 0, err
	}

	var get models.Response
	json.NewDecoder(resp.Body).Decode(&get)
	result, err := strconv.ParseFloat(fmt.Sprintf("%.2f", euro*get.Rates[currency]), 32)
	if err != nil {
		return 0, err
	}

	return float32(result), nil
}

func ParseTime(value string, t *testing.T) time.Time {
	timeAt, err := time.Parse(time.DateTime, value)
	if err != nil {
		t.Error(err)
	}

	return timeAt
}

func Float2String(xF float32) string {
	if math.Trunc(float64(xF)) == float64(xF) {
		return fmt.Sprintf("%.0f", xF)
	} else {
		return fmt.Sprintf("%.2f", xF)
	}
}
