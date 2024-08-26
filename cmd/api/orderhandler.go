package main

import (
	"fmt"
	"net/http"
)

func (app *application) getCreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Get all orders")
}
