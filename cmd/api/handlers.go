package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/esousacosta/managementsystem/internal/data"
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

		err = writeJson(w, http.StatusOK, envelope{"books": books}, nil)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	case r.Method == http.MethodPost:
		var input struct {
			Id        int       `json:"-"`
			CreatedAt time.Time `json:"-"`
			Name      string    `json:"name"`
			Price     float32   `json:"price"`
			Stock     int64     `json:"stock"`
			Reference string    `json:"reference"`
			BarCode   string    `json:"barcode"`
		}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&input)
		if err != nil {
			app.logger.Printf("decoding error --> %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		part := data.Part{
			Name:      input.Name,
			Price:     input.Price,
			Stock:     input.Stock,
			Reference: input.Reference,
			BarCode:   input.BarCode,
		}

		err = app.model.Parts.Insert(&part)
		if err != nil {
			app.logger.Printf("insertion error --> %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Create new book: %v", input)
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
