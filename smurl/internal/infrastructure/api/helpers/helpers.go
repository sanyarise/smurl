package helpers

import (
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

//Вспомогательная функция для проверки валидности длинного урл
func CheckURL(longURL string) bool {
	url, err := url.Parse(longURL)
	return err == nil && url.Scheme != "" && url.Host != ""
}

//Вспомогательная функция для генерации случайной строки,
//использующейся для уменьшенного и админского урл
func RandString(l *zap.Logger) string {
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		l.Error("error on rand.read(buf)",
			zap.Error(err))
		return ""
	}
	return fmt.Sprintf("%x", buf)
}

//Вспомогательная функция для увеличения счетчика
//переходов по уменьшенному урл
func CountUses(s string, l *zap.Logger) (string, error) {
	l.Debug("CountUses")
	count, err := strconv.Atoi(s)
	if err != nil {
		l.Error("",
			zap.Error(err))
		return "", err
	}
	l.Debug("Count",
		zap.Int(":", count))
	count++
	l.Debug("Count upd",
		zap.Int(":", count))
	s = strconv.Itoa(count)
	return s, nil
}

//Вспомогательная функция, позволяющая прочитать
//информацию об ip адресе из реквеста
func GetIP(r *http.Request) string {
	rip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(rip)
	if netIP != nil {
		return rip
	}
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown ip"
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip
	}
	return "unknown ip"
}
