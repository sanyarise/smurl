package helpers

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Helper interface {
	CheckURL(longURL string) bool
	GetIP(r *http.Request) string
	RandString() string
}
type Helpers struct {
	logger *zap.Logger
}

func NewHelpers(logger *zap.Logger) Helper {
	return &Helpers{logger: logger}
}

const alphabet = "ynAJfoSgdXHB5VasEMtcbPCr1uNZ4LG723ehWkvwYR6KpxjTm8iQUFqz9D"

var alphabetLen = uint32(len(alphabet))

// CheckURL check the validity of a long url
func (helpers *Helpers) CheckURL(longURL string) bool {
	helpers.logger.Debug("Enter in func CheckURL()")
	url, err := url.ParseRequestURI(longURL)
	return err == nil && url.Scheme != "" && url.Host != "" && strings.Contains(url.Host, ".")
}

// RandString generate a random string,
// used for minified and admin url
func (helpers *Helpers) RandString() string {
	helpers.logger.Debug("Enter in func RandString()")
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

// GetIP read information about
// the ip address from the request
func (helpers *Helpers) GetIP(r *http.Request) string {
	helpers.logger.Debug("Enter in func GetIP()")
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
