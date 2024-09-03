package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/auth", app.postCreateAuthenticationHandler)
	mux.HandleFunc("/v1/auth/", app.postRequestAuthenticationHandler)
	mux.HandleFunc("/v1/parts", app.validateSession(app.getCreatePartsHandler))
	mux.HandleFunc("/v1/parts/", app.getUpdateDeletePartsHandler)
	mux.HandleFunc("/v1/orders", app.getCreateOrdersHandler)
	mux.HandleFunc("/v1/orders/search", app.getFilteredOrdersHandler)
	mux.HandleFunc("/v1/orders/", app.getUpdateDeleteOrdersHandler)
	mux.HandleFunc("/v1/clients", app.validateSession(app.getCreateClientsHandler))
	mux.HandleFunc("/v1/clients/", app.getUpdateDeleteClientsHandler)
	return mux
}
