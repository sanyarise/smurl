package helpers

import (
	"crypto/rand"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// CheckURL check the validity of a long url
func CheckURL(longURL string) bool {
	url, err := url.Parse(longURL)
	return err == nil && url.Scheme != "" && url.Host != ""
}

// RandString generate a random string,
// used for minified and admin url
func RandString(l *zap.Logger) string {
	buf := make([]byte, 4)
	num, err := rand.Read(buf)
	if err != nil {
		l.Error("error on rand.Read(buf)",
			zap.Error(err))
		return ""
	}
	res := Base62Encode(uint64(num))
	return res
}

func Base62Encode(number uint64) string {
	const (
		alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	)
	length := len(alphabet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(alphabet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

// CountUses increment the counter
// transitions on a reduced url
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

// GetIP read information about
// the ip address from the request
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
