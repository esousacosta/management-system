package main

import (
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) validateSession(finalHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, "authentication failed. Please log in", http.StatusUnauthorized)
			return
		}

		log.Printf("Received cookis during validation: %v", cookie.Value)

		_, _ = jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(jwt.RegisteredClaims)
			if !ok {
				app.logger.Printf("[%s] ERROR - Invalid JWT claims", shared.GetCallerInfo())
				http.Error(w, "authentication failed. Please log in", http.StatusUnauthorized)
			}
			app.logger.Printf("cookies: %v", cookie)
			app.logger.Printf("user email: %s", claims.ID)
			return nil, nil
		})
		finalHandler(w, r)
	}
}
