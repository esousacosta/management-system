package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/esousacosta/managementsystem/cmd/shared"
)

func (app *application) getCreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		orders, err := app.model.Orders.GetAll()
		if err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := writeJson(w, http.StatusOK, envelope{"orders": orders}, nil); err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		order, err := readOrderJson(w, r, app.logger)
		if err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = app.model.Orders.Insert(order)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		header := make(http.Header)
		header.Set("Location", "/v1/orders/"+strconv.Itoa(order.ID))
		if err := writeJson(w, http.StatusCreated, envelope{"order": order}, header); err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) getUpdateDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := shared.GetUniqueIdentifierFromUrl("/v1/orders/", r)
		id, err := strconv.ParseInt(*idStr, 10, 64)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		order, err := app.model.Orders.Get(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := writeJson(w, http.StatusOK, envelope{"order": *order}, nil); err != nil {
			app.logger.Printf("error sending the response: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		return
	case http.MethodPut:
		fmt.Fprintln(w, "Update a single order")
		return
	case http.MethodDelete:
		fmt.Fprintln(w, "Delete a single order")
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}
