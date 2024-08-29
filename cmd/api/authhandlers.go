package main

import "net/http"

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
	hashedPassword := hashPassword()
}
