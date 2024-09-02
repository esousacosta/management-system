package main

import (
	"net/http"
)

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/login", app.loginHandler)
	mux.HandleFunc("/unauthorized", app.getUnauthorizedHandler)
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/order/create", app.createOrder)
	mux.HandleFunc("/orders/search", app.filteredOrdersView)
	mux.HandleFunc("/parts", app.partsView)
	mux.HandleFunc("/part/view/", app.viewPart)
	mux.HandleFunc("/part/create", app.createPart)

	return mux
}
