package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "server alive\n")
}

func (app *application) getCreatePartsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		parts, err := app.model.Parts.GetAll()
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = writeJson(w, http.StatusOK, envelope{"parts": parts}, nil)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case r.Method == http.MethodPost:

		part, err := readJson(w, r, app.logger)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = app.model.Parts.Insert(part)
		if err != nil {
			app.logger.Printf("insertion error --> %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/parts/%d", part.Id))

		if err := writeJson(w, http.StatusCreated, envelope{"part": part}, headers); err != nil {
			app.logger.Printf("db writing error --> %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getUpdateDeletePartsHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		fmt.Fprintf(w, "Get a single part")
	case r.Method == http.MethodPut:
		fmt.Fprintf(w, "Update a part")
	case r.Method == http.MethodDelete:
		fmt.Fprintf(w, "Delete a part")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
