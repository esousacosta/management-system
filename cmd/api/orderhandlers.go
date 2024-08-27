package main

import (
	"fmt"
	"net/http"
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

		app.model.Orders.Insert(order)
		fmt.Fprintf(w, "%+v", order)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) getUpdateDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Get a single order")
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
