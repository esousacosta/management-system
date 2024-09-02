package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/golang-jwt/jwt/v4"
)

type Key int

const myKey Key = 0

func (app *application) validateSession(finalHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth")
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, "authentication failed. Please log in", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, fmt.Errorf("[%s] ERROR - Invalid JWT claims", shared.GetCallerInfo())
			}
			app.logger.Printf("cookies: %v", cookie)
			app.logger.Printf("user email: %s", claims["id"])
			userEmail, ok := claims["id"].(string)
			if !ok {
				return nil, fmt.Errorf("[%s] ERROR - Invalid User Email", shared.GetCallerInfo())
			}
			app.logger.Print("zzzzzzzzzzzzz ALL TOP zzzzzzzzzzzzz")
			userAuth, err := app.model.UsersAuth.GetUserAuth(userEmail)
			if err != nil {
				return nil, fmt.Errorf("[%s] ERROR - couldn't retrieve user", shared.GetCallerInfo())
			}
			app.logger.Print("XXXXXXXXX ALL GOOD XXXXXXXXX")
			return []byte(userAuth.JwtSecret), nil
		})
		if err != nil {
			app.logger.Print("zzzzzzzzzzzzz ALL BAD zzzzzzzzzzzzz")
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, "authentication failed. Please log in", http.StatusUnauthorized)
			return
		}

		app.logger.Print("KKKKKKKKKK ALL GOOD KKKKKKKKKK")
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			app.logger.Printf("[%s] ERROR - unexpected signing method", shared.GetCallerInfo())
			http.Error(w, "authentication failed. Please log in", http.StatusUnauthorized)
			return
		}

		app.logger.Print("HHHHHHHHH ALL GOOD HHHHHHHHH")
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), myKey, claims)
			app.logger.Print("YYYYYYYYYYYYY ALL STILL GOOD YYYYYYYYYYYYY")
			finalHandler(w, r.WithContext(ctx))
		} else {
			app.logger.Printf("[%s] ERROR - invalid token or claims", shared.GetCallerInfo())
			http.Error(w, "Unauthorized - please log in", http.StatusUnauthorized)
			return
		}
	}
}
