package shared

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecretSize    int    = 32
	webServerEnpoint string = "localhost:4000"
)

func GetViewsFuncMap() template.FuncMap {
	return template.FuncMap{
		"timeFormatting": func() string { return time.DateTime },
		"seq": func(start, end int) []int {
			seq := []int{}
			for i := start; i <= end; i++ {
				seq = append(seq, i)
			}
			return seq
		},
	}
}

func GetUniqueIdentifierFromUrl(url string, r *http.Request) *string {
	ref := r.URL.Path[len(url):]
	return &ref
}

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 14)
	return string(hashedPasswordBytes), err
}

func GenerateUnsignedJwtToken(userEmail string) (*jwt.Token, error) {
	// 30m expiration for non-sensitive applications - OWASP
	tokenExpirationTime := time.Now().Add(time.Minute * 30)

	claims := jwt.MapClaims{
		"id":  userEmail,
		"exp": jwt.NewNumericDate(tokenExpirationTime),
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims), nil

}

func GenerateUserRandomJwtSecret() (string, error) {
	secret := make([]byte, jwtSecretSize)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(secret), nil
}

func GetCallerInfo() string {
	pc, file, line, ok := runtime.Caller(1)
	if ok {
		lineNumber := strconv.Itoa(line)
		return file + " (line #" + lineNumber + ") - " + runtime.FuncForPC(pc).Name()
	}
	return ""
}

func GetCertPool() *x509.CertPool {
	certPool := x509.NewCertPool()
	serverCert, err := os.ReadFile("./cert/domain.crt")
	if err != nil {
		log.Fatalf("[%s - ERROR] %s", GetCallerInfo(), err)
	}

	certPool.AppendCertsFromPEM(serverCert)
	return certPool
}
