package main

import (
	"net/http"
	"strconv"

	"github.com/esousacosta/managementsystem/cmd/shared"
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
	default:
		app.logger.Printf("[%s] ERROR - unexpected request method %s", shared.GetCallerInfo(), r.Method)
		// case http.MethodPost:
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
