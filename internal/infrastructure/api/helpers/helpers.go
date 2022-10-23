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
	buf := make([]byte, 4)
	_, err := rand.Read(buf)
	if err != nil {
		l.Error("error on rand.read(buf)",
			zap.Error(err))
		return ""
	}
	l.Debug("Random String generate success")
	return fmt.Sprintf("%x", buf)
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
