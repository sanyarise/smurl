package helpers

import (
	"bytes"
	"log"
	"net/http"
	"testing"

	"go.uber.org/zap"
)

func TestCheckURL(t *testing.T) {
	var tests = []struct {
		url  string
		wait bool
	}{
		{"http://ok.ru", true},
		{"https://google.com", true},
		{"vk.com", false},
		{"", false},
		{"afsdfasdfasd", false},
	}
	for i, e := range tests {
		res := CheckURL(tests[i].url, zap.L())
		if res != e.wait {
			t.Errorf("CheckURL(%s) =  %v, wait %v", tests[i].url, res, tests[i].wait)
		}
	}
}

func TestCountUses(t *testing.T) {
	l := zap.L()
	var tests = []struct {
		count    string
		countNew string
	}{
		{"0", "1"},
		{"1", "2"},
		{"2", "3"},
		{"10000", "10001"},
		{"59382920", "59382921"},
		{"afd", ""},
		{"", ""},
	}
	for i, e := range tests {
		res, err := CountUses(e.count, l)
		if err != nil {
			log.Println(err)
		}
		if res != tests[i].countNew {
			t.Errorf("CountUses(%s) = %s, wait %s", tests[i].count, res, tests[i].countNew)
		}
	}
}

func TestGetIP(t *testing.T) {
	jsonStr := []byte("")
	ip := "0.0.0.0"
	r, _ := http.NewRequest("GET", "http://vk.com", bytes.NewBuffer(jsonStr))
	r.Header.Set("X-REAL-IP", ip)
	res := GetIP(r, zap.L())
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.Header.Set("X-FORWARDED-FOR", ip)
	res = GetIP(r, zap.L())
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.RemoteAddr = "[::1]:80"
	r.Header.Set("X-REAL-IP", "")
	r.Header.Set("X-FORWARDED-FOR", "")
	ip = "::1"
	res = GetIP(r, zap.L())
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
	r.RemoteAddr = ""
	r.Header.Set("X-REAL-IP", "")
	r.Header.Set("X-FORWARDED-FOR", "")
	ip = "unknown ip"
	res = GetIP(r, zap.L())
	if res != ip {
		t.Errorf("error GetIP: wait %s, got %s", ip, res)
	}
}
