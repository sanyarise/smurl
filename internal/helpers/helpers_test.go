package helpers

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCheckURL(t *testing.T) {
	logger := zap.L()
	helpers := NewHelpers(logger)

	var tests = []struct {
		url    string
		expect bool
	}{
		{"http://ok.ru", true},
		{"https://google.com", true},
		{"vk.com", false},
		{"", false},
		{"afsdfasdfasd", false},
		{"http://example.ru/sadfasdfasdfasfasfasfdasdf?Name=safdasf&id=2", true},
		{"NaN", false},
		{"httP:/mail.ru", false},
		{"HTTP://MAIL.RU", true},
		{"http://mail", false},
		{"https://mail.12", true},
	}
	for i, test := range tests {
		res := helpers.CheckURL(tests[i].url)
		assert.Equal(t, test.expect, res)
	}
}

func TestGetIP(t *testing.T) {
	logger := zap.L()
	helpers := NewHelpers(logger)

	jsonStr := []byte("")
	ip := "0.0.0.0"
	r, _ := http.NewRequest("GET", "http://vk.com", bytes.NewBuffer(jsonStr))
	r.Header.Set("X-REAL-IP", ip)
	res := helpers.GetIP(r)
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.Header.Set("X-FORWARDED-FOR", ip)
	res = helpers.GetIP(r)
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.RemoteAddr = "[::1]:80"
	r.Header.Set("X-REAL-IP", "")
	r.Header.Set("X-FORWARDED-FOR", "")
	ip = "::1"
	res = helpers.GetIP(r)
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.RemoteAddr = ""
	r.Header.Set("X-REAL-IP", "")
	r.Header.Set("X-FORWARDED-FOR", "")
	ip = "unknown ip"
	res = helpers.GetIP(r)
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
}
