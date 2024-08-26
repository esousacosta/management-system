package main

import (
	"fmt"
	"net/http"
)

func (app *application) getCreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Get all orders")
		return
	case http.MethodPost:
		fmt.Fprintln(w, "Create an order")
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
