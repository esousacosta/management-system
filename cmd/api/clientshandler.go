package main

import (
	"net/http"
	"strconv"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) getCreateClientsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		clients, err := app.model.Clients.GetAll()
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = writeJson(w, http.StatusOK, envelope{"clients": clients}, nil)
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		app.createClient(w, r)
	default:
		app.logger.Printf("[%s] ERROR - unexpected request method %s", shared.GetCallerInfo(), r.Method)
		return
	}
}

func (app *application) createClient(w http.ResponseWriter, r *http.Request) {
	readClient, err := readClientJson(w, r, app.logger)
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value(myKey).(jwt.MapClaims)
	if !ok {
		app.logger.Printf("[%s] ERROR - invalid JWT token in the request", shared.GetCallerInfo())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Golang decoded numbers to float64 by default from its unmarshaller
	// TODO: change the claims type from MapClaim to something else
	userId, ok := claims["id"].(float64)
	if !ok {
		app.logger.Printf("[%s] ERROR - invalid <<id>> present in the request's JWT token", shared.GetCallerInfo())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = app.model.Clients.Insert(readClient, int(userId))
	if err != nil {
		app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", "/v1/clients/"+strconv.Itoa(readClient.Id))

	if err := writeJson(w, http.StatusCreated, envelope{"client": *readClient}, headers); err != nil {
		app.logger.Printf("[%s] response writing error --> %v", shared.GetCallerInfo(), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) getUpdateDeleteClientsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		clientIdStr := shared.GetUniqueIdentifierFromUrl("/v1/clients/", r)
		clientId, err := strconv.Atoi(*clientIdStr)
		if err != nil {
			app.logger.Printf("[%s] ERROR parsing client ID - %v", shared.GetCallerInfo(), err)
			http.Error(w, "Invalid client ID requested", http.StatusBadRequest)
			return
		}

		client, err := app.model.Clients.GetClientById(clientId)
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = writeJson(w, http.StatusOK, envelope{"client": *client}, nil)
		if err != nil {
			app.logger.Printf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		app.logger.Printf("[%s] ERROR - unexpected request method %s", shared.GetCallerInfo(), r.Method)
		// case http.MethodPost:
	}
}
