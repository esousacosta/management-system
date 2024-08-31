package shared

import (
	"net/http"
	"runtime"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

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
