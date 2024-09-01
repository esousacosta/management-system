package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/esousacosta/managementsystem/cmd/shared"
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
		part, err := readPartJson(w, r, app.logger)
		if err != nil {
			app.logger.Printf("request decoding error --> %v", err)
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
		headers.Set("Location", "/v1/parts/"+strconv.Itoa(part.Id))

		if err := writeJson(w, http.StatusCreated, envelope{"part": part}, headers); err != nil {
			app.logger.Printf("response writing error --> %v", err)
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
		app.deletePart(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (app *application) getPart(w http.ResponseWriter, r *http.Request) {
	ref := shared.GetUniqueIdentifierFromUrl("/v1/parts/", r)
	part, err := app.model.Parts.GetByRef(*ref)
	if err != nil {
		app.logger.Printf("[%s] ERROR - %s", shared.GetCallerInfo(), err.Error())
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
	ref := shared.GetUniqueIdentifierFromUrl("/v1/parts/", r)
	part, err := app.model.Parts.GetByRef(*ref)
	if err != nil {
		app.logger.Printf("[%s] ERROR - %s", shared.GetCallerInfo(), err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	receivedPart, err := readPartJson(w, r, app.logger)
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

	if receivedPart.Barcode != nil {
		part.Barcode = *receivedPart.Barcode
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
	app.logger.Printf("Part with reference %s updated", *ref)
}

func (app *application) deletePart(w http.ResponseWriter, r *http.Request) {
	ref := shared.GetUniqueIdentifierFromUrl("/v1/parts/", r)
	err := app.model.Parts.Delete(*ref)
	if err != nil {
		switch {
		case errors.Is(err, fmt.Errorf("no part found with ref %s", *ref)):
			app.logger.Printf("[ERROR] - %v", err)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		default:
			app.logger.Printf("[ERROR] - %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	if err := writeJson(w, http.StatusOK, envelope{"message": fmt.Sprintf("part with reference %s deleted", *ref)}, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	app.logger.Printf("Part with reference %s deleted", *ref)
}
