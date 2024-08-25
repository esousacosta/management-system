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
		ref := r.URL.Query().Get("ref")
		app.logger.Printf("Retrieved ref: %s", ref)
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
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
		app.getPart(w, r)
	case r.Method == http.MethodPut:
		app.updatePart(w, r)
	case r.Method == http.MethodDelete:
		fmt.Fprintf(w, "Delete a part")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getPart(w http.ResponseWriter, r *http.Request) {
	ref := getPartReferenceFromUrl("/v1/parts/", r)
	part, err := app.model.Parts.GetByRef(*ref)
	if err != nil {
		app.logger.Printf("[ERROR] - %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := writeJson(w, http.StatusOK, envelope{"part": part}, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) updatePart(w http.ResponseWriter, r *http.Request) {
	ref := getPartReferenceFromUrl("/v1/parts/", r)
	part, err := app.model.Parts.GetByRef(*ref)
	if err != nil {
		app.logger.Printf("[ERROR] - %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	receivedPart, err := readJson(w, r, app.logger)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if receivedPart.Name != nil {
		part.Name = *receivedPart.Name
	}

	if receivedPart.Price != nil {
		part.Price = *receivedPart.Price
	}

	if receivedPart.Reference != nil {
		part.Reference = *receivedPart.Reference
	}

	if receivedPart.BarCode != nil {
		part.BarCode = *receivedPart.BarCode
	}

	err = app.model.Parts.Update(part)
	if err != nil {
		app.logger.Printf("[ERROR] - %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := writeJson(w, http.StatusOK, envelope{"updated part": part}, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.logger.Printf("Part with reference %s successfully updated", *ref)
}
