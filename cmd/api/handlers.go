package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server alive\n")
}

func (app *application) getCreatePartsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		books, err := app.model.Parts.GetAll()
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		json, err := json.MarshalIndent(map[string]any{"books": books}, "", "\t")
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(json)
		// fmt.Fprintf(w, "books: %v", map[string]any{"books": books})
	case r.Method == http.MethodPost:
		fmt.Fprintf(w, "Create new book")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeletePartsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		fmt.Fprintf(w, "Get a single book")
	case r.Method == http.MethodPut:
		fmt.Fprintf(w, "Update a book")
	case r.Method == http.MethodDelete:
		fmt.Fprintf(w, "Delete a book")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
