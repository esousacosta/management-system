package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/esousacosta/managementsystem/internal/data"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/orders.html",
		"./ui/html/partials/nav.html",
	}

	orders, err := app.managSysModel.GetAllOrders()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"timeFormatting": func() string { return time.DateTime },
	}

	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", orders)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) filteredOrdersView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	clientId := r.FormValue("clientid")
	if clientId == "" {
		log.Print("error reading client ID from the form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/orders.html",
		"./ui/html/partials/nav.html",
	}

	orders, err := app.managSysModel.GetOrdersByClientId(clientId)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"timeFormatting": func() string { return time.DateTime },
	}

	ts, err := template.New("base").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", orders)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) createOrder(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.createOrderForm(w, r)
	case http.MethodPost:
		app.createOrderProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) createOrderForm(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/createOrder.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) createOrderProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("create order form parsing error --> %v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	clientId := r.PostFormValue("client_id")
	if clientId == "" {
		log.Printf("client_id form reading error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	getFormArrayValueAsSlice := func(r *http.Request, cannotBeNull bool, key string) ([]string, error) {
		formValue := r.PostFormValue(key)
		if cannotBeNull && formValue == "" {
			log.Printf("%s form reading error", key)
			return nil, fmt.Errorf("empty form field %s", key)
		}
		splitStrings := strings.Split(formValue, ",")
		for i, splitStr := range splitStrings {
			splitStrings[i] = strings.TrimSpace(splitStr)
		}

		return splitStrings, nil
	}

	services, err := getFormArrayValueAsSlice(r, false, "services")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	partsRefs, err := getFormArrayValueAsSlice(r, false, "parts_refs")
	if err != nil {
		log.Printf("parts_refs form reading error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	comment := r.PostFormValue("comment")

	totalFloat, err := strconv.ParseFloat(r.PostFormValue("total"), 32)
	if err != nil || totalFloat < 0 {
		log.Printf("total form reading error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	total := float32(totalFloat)

	order := &data.Order{
		ClientId:  clientId,
		Services:  services,
		PartsRefs: partsRefs,
		Comment:   comment,
		Total:     total,
	}

	errorCode := app.managSysModel.PostOrder(order)
	if errorCode != http.StatusCreated {
		http.Error(w, http.StatusText(int(errorCode)), int(errorCode))
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) partsView(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/pages/parts.html",
		"./ui/html/partials/nav.html",
	}

	parts, err := app.managSysModel.GetAllParts()
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
	ref := shared.GetUniqueIdentifierFromUrl("/part/view/", r)
	fmt.Println(*ref)
	part, err := app.managSysModel.GetPart(*ref)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/view.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", *part)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) createPart(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		app.createPartForm(w, r)
	case http.MethodPost:
		app.createPartProcess(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) createPartForm(w http.ResponseWriter, _ *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/createPart.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (app *application) createPartProcess(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("create part form parsing error --> %v", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		log.Print("error reading name from the form")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	priceFloat, err := strconv.ParseFloat(r.PostFormValue("price"), 32)
	if err != nil || priceFloat < 0 {
		log.Printf("error reading the price from the form")
		if err != nil {
			log.Print(err.Error())
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	price := float32(priceFloat)

	stock, err := strconv.ParseInt(r.PostForm.Get("stock"), 10, 64)
	if err != nil || stock < 0 {
		log.Printf("error reading the stock from the form")
		if err != nil {
			log.Print(err.Error())
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	reference := r.PostFormValue("reference")
	if reference == "" {
		log.Printf("the part reference cannot be empty")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	barcode := r.PostFormValue("barcode")

	part := &data.Part{
		Name:      name,
		Price:     price,
		Stock:     stock,
		Reference: reference,
		Barcode:   barcode,
	}

	errorCode := app.managSysModel.PostPart(part)
	if errorCode != http.StatusCreated {
		http.Error(w, http.StatusText(int(errorCode)), int(errorCode))
		return
	}

	http.Redirect(w, r, "/parts", http.StatusSeeOther)
}
