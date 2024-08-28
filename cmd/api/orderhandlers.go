package main

import (
	"net/http"
	"strconv"

	"github.com/esousacosta/managementsystem/cmd/shared"
)

func (app *application) getCreateOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		orders, err := app.model.Orders.GetAll()
		if err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := writeJson(w, http.StatusOK, envelope{"orders": orders}, nil); err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	case http.MethodPost:
		order, err := readOrderJson(w, r, app.logger)
		if err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = app.model.Orders.Insert(order)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		header := make(http.Header)
		header.Set("Location", "/v1/orders/"+strconv.Itoa(order.ID))
		if err := writeJson(w, http.StatusCreated, envelope{"order": order}, header); err != nil {
			app.logger.Print(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (app *application) getUpdateDeleteOrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		idStr := shared.GetUniqueIdentifierFromUrl("/v1/orders/", r)
		id, err := strconv.ParseInt(*idStr, 10, 64)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		order, err := app.model.Orders.Get(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := writeJson(w, http.StatusOK, envelope{"order": *order}, nil); err != nil {
			app.logger.Printf("error sending the response: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		return
	case http.MethodPut:
		idStr := shared.GetUniqueIdentifierFromUrl("/v1/orders/", r)
		id, err := strconv.ParseInt(*idStr, 10, 64)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		order, err := app.model.Orders.Get(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		receivedOrder, err := readOrderJson(w, r, app.logger)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if receivedOrder.ClientId != nil {
			order.ClientId = *receivedOrder.ClientId
		}

		if receivedOrder.Services != nil {
			order.Services = *receivedOrder.Services
		}

		if receivedOrder.PartsRefs != nil {
			order.PartsRefs = *receivedOrder.PartsRefs
		}

		if receivedOrder.Comment != nil {
			order.Comment = *receivedOrder.Comment
		}

		if receivedOrder.Total != nil {
			order.Total = *receivedOrder.Total
		}

		err = app.model.Orders.Update(order)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = writeJson(w, http.StatusOK, envelope{"updated order": order}, nil)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		app.logger.Printf("Order with reference %s updated", *idStr)
	case http.MethodDelete:
		idStr := shared.GetUniqueIdentifierFromUrl("/v1/orders/", r)
		id, err := strconv.ParseInt(*idStr, 10, 64)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = app.model.Orders.Delete(id)
		if err != nil {
			app.logger.Print(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = writeJson(w, http.StatusOK, envelope{"order with id %s deleted": id}, nil)
		if err != nil {
			app.logger.Printf("error writing response from server: %v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		app.logger.Printf("Deleted order with id %d", id)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
