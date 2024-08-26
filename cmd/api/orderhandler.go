package main

import (
	"fmt"
	"net/http"
)

func (app *application) getCreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Get all orders")
}

func (app *application) getUpdateDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "Get a single order")
	case http.MethodPost:
		fmt.Fprintln(w, "Insert a single order")
	case http.MethodDelete:
		fmt.Fprintln(w, "Delete a single order")
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Get all orders")
}
