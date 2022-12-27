package helpers

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

// CheckURL check the validity of a long url
func CheckURL(longURL string, l *zap.Logger) bool {
	l.Debug("Enter in func CheckURL()")
	url, err := url.Parse(longURL)
	return err == nil && url.Scheme != "" && url.Host != ""
}

// RandString generate a random string,
// used for minified and admin url
func RandString(l *zap.Logger) string {
	l.Debug("Enter in func RandString()")
	var (
		digits  []uint32
		num     = uuid.New().ID()
		builder strings.Builder
	)
	for num > 0 {
		digits = append(digits, num%alphabetLen)
		num /= alphabetLen
	}
	digits = Reverse(digits)

	for _, digit := range digits {
		builder.WriteString(string(alphabet[digit]))
	}
	return builder.String()
}

// CountUses increment the counter
// transitions on a reduced url
func CountUses(s string, l *zap.Logger) (string, error) {
	l.Debug("Enter in func CountUses()")
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
func GetIP(r *http.Request, l *zap.Logger) string {
	l.Debug("Enter in func GetIP()")
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

// Revers reverse slice of uint32
func Reverse(param []uint32) []uint32 {
	result := make([]uint32, len(param))
	counter := 0
	for i := len(param) - 1; i >= 0; i-- {
		result[counter] = param[i]
		counter++
	}
	return result
}
