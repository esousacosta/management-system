package shared

import (
	"crypto/x509"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
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
