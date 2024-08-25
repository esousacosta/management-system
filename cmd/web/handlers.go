package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Simplex management system!")
}

func (app *application) viewPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "View of a single part")
}

func (app *application) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creation form for a single part")
}
