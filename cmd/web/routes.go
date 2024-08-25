package main

import (
	"fmt"
	"net/http"
	"os"
)

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()

	fmt.Println(os.Getwd())
	fileServer := http.FileServer(http.Dir("./../../ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/part/view", app.viewPart)
	mux.HandleFunc("/part/create", app.createPart)

	return mux
}
