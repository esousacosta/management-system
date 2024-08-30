package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/auth", app.postCreateAuthenticationHandler)
	mux.HandleFunc("/v1/auth/", app.postRequestAuthenticationHandler)
	mux.HandleFunc("/v1/parts", app.getCreatePartsHandler)
	mux.HandleFunc("/v1/parts/", app.getUpdateDeletePartsHandler)
	mux.HandleFunc("/v1/orders", app.getCreateOrdersHandler)
	mux.HandleFunc("/v1/orders/search", app.getFilteredOrdersHandler)
	mux.HandleFunc("/v1/orders/", app.getUpdateDeleteOrdersHandler)
	return mux
}
