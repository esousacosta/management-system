package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/part/view", app.viewPart)
	mux.HandleFunc("/part/create", app.createPart)

	return mux
}
