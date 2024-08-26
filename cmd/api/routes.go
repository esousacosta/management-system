package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthcheck", app.healthcheck)
	mux.HandleFunc("/v1/parts", app.getCreatePartsHandler)
	mux.HandleFunc("/v1/parts/", app.getUpdateDeletePartsHandler)
	mux.HandleFunc("/v1/orders", app.getCreateOrdersHandler)
	return mux
}
