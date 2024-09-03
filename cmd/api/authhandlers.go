package main

import (
	"net/http"
	"time"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/esousacosta/managementsystem/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) postCreateAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.createUserAuth(w, r)
	default:
		http.Error(w, "invalid requested method", http.StatusMethodNotAllowed)
	}
}

func (app *application) createUserAuth(w http.ResponseWriter, r *http.Request) {
	userAuth, err := readJson[data.ReadUserAuth](w, r, app.logger)
	if err != nil {
		app.logger.Print("[ERROR] couldn't parse the received request's body")
		http.Error(w, "couldn't parse the received form", http.StatusBadRequest)
		return
	}

	hashedPassword, err := shared.HashPassword(*userAuth.Password)
	if err != nil {
		app.logger.Print("error hashing the received password")
		http.Error(w, "Error parsing the user credentials", http.StatusInternalServerError)
		return
	}
	userAuth.Password = &hashedPassword

	userJwtSecret, err := shared.GenerateUserRandomJwtSecret()
	if userJwtSecret == "" {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, "error processing user creation request. Try again later.", http.StatusInternalServerError)
	}
	userAuth.JwtSecret = &userJwtSecret

	err = app.model.UsersAuth.InsertUser(userAuth)
	if err != nil {
		app.logger.Print(err.Error())
		http.Error(w, "error processing user creation request", http.StatusBadRequest)
		return
	}

	if err := writeJson(w, http.StatusCreated, envelope{"response": "user <" + *userAuth.Email + "> created successfully"}, nil); err != nil {
		app.logger.Printf("response writing error --> %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) postRequestAuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		app.processUserAuth(w, r)
	default:
		http.Error(w, "invalid requested method", http.StatusMethodNotAllowed)
	}
}

func (app *application) processUserAuth(w http.ResponseWriter, r *http.Request) {
	userAuth, err := readJson[data.ReadUserAuth](w, r, app.logger)
	if err != nil {
		app.logger.Print("[ERROR] couldn't parse the received request's body")
		http.Error(w, "couldn't parse the received form", http.StatusBadRequest)
		return
	}

	user, err := app.model.UsersAuth.GetUserAuth(*userAuth.Email)
	if err != nil {
		http.Error(w, "incorrect user credentials", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*userAuth.Password))
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		err := writeJson(w, http.StatusUnauthorized, envelope{"authentication": false}, nil)
		if err != nil {
			http.Error(w, "error writing response", http.StatusInternalServerError)
			return
		}
		return
	}

	// Set JWT Token on the response cookie
	unsignedToken, err := shared.GenerateUnsignedJwtToken(user.Email, user.Id)
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, "error generating authentication token", http.StatusInternalServerError)
		return
	}

	// the key needs to be of type []byte
	signedToken, err := unsignedToken.SignedString([]byte(user.JwtSecret))
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, "error generating authentication token", http.StatusInternalServerError)
		return
	}

	cookieExpirationTime := time.Now().Add(time.Minute * 30)
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    signedToken,
		Expires:  cookieExpirationTime,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Domain:   "localhost",
	}

	headers := make(http.Header)
	// CORS setup
	headers["Access-Control-Allow-Origin"] = []string{"https://localhost:3000/"}
	headers["Access-Control-Allow-Credentials"] = []string{"true"}

	http.SetCookie(w, cookie)

	err = writeJson(w, http.StatusOK, envelope{"authenticated": true}, headers)
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, "error writing response", http.StatusInternalServerError)
		return
	}
}
