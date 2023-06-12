package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	date := time.Now()
	dateStr := date.Format(time.DateTime)
	res, err := time.Parse(time.DateTime, dateStr)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, res, ParseTime(dateStr, t))
}
