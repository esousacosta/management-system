package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/cmd/shared"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"../../ui/html/base.html",
		"../../ui/html/pages/home.html",
		"../../ui/html/partials/nav.html",
	}

	parts, err := app.managSysModel.GetAll()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", parts)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) viewPart(w http.ResponseWriter, r *http.Request) {
	ref := shared.GetPartReferenceFromUrl("/part/view/", r)
	fmt.Println(*ref)
	// app.managSysModel.getPart()
	fmt.Fprintf(w, "View of a single part")
}

func (app *application) createPart(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Creation form for a single part")
}
