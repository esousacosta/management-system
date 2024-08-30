package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (app *application) getPostAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.processUserAuth(w, r)
	default:
		http.Error(w, "invalid requested method", http.StatusMethodNotAllowed)
	}
}

func (app *application) processUserAuth(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.logger.Print("[ERROR] couldn't parse the received form")
		http.Error(w, "couldn't parse the received form", http.StatusBadRequest)
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	hashedPassword, err := hashPassword(password)
	if err != nil {
		app.logger.Print("error hashing the received password")
		http.Error(w, "Error parsing the user credentials", http.StatusInternalServerError)
		return
	}

	user, err := app.model.UsersAuth.GetUser(email)
	if err != nil {
		http.Error(w, "incorrect user credentials", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "incorrect user credentials", http.StatusUnauthorized)
		return
	}

	err = writeJson(w, http.StatusOK, envelope{"authorization": true}, nil)
	if err != nil {
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
}

func hashPassword(password string) (string, error) {
	passwordBytes := []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 14)
	return string(hashedPasswordBytes), err
}
