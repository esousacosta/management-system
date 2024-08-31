package appmodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/esousacosta/managementsystem/internal/data"
)

type errorCode int

type ManagementSystemModel struct {
	PartsEndpoint  string
	OrdersEndpoint string
	AuthEndpoint   string
}

type PartsResponse struct {
	Parts []data.Part `json:"parts"`
}

type OrdersResponse struct {
	Orders []data.Order `json:"orders"`
}

type PartResponse struct {
	Parts data.Part `json:"part"`
}

type AuthResponse struct {
	Authorized bool `json:"authenticated"`
}

func NewManagementSystemModel(ordersEndpoint string, partsEndpoint string, authEndpoint string) ManagementSystemModel {
	return ManagementSystemModel{
		PartsEndpoint:  partsEndpoint,
		OrdersEndpoint: ordersEndpoint,
		AuthEndpoint:   authEndpoint,
	}
}

func (managSysModel *ManagementSystemModel) GetAllParts() (*[]data.Part, error) {
	resp, err := http.Get(managSysModel.PartsEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var partsResp PartsResponse

	err = json.Unmarshal(data, &partsResp)
	if err != nil {
		return nil, err
	}

	return &partsResp.Parts, nil
}

func (managSysModel *ManagementSystemModel) GetPart(partRef string) (*data.Part, error) {
	resp, err := http.Get(managSysModel.PartsEndpoint + "/" + partRef)
	if err != nil {
		return nil, fmt.Errorf("part with reference %s not found", partRef)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status received: %s", resp.Status)
	}

	var partResponse PartResponse

	err = json.Unmarshal(data, &partResponse)
	if err != nil {
		return nil, err
	}

	return &partResponse.Parts, nil
}

func (managSysMoel *ManagementSystemModel) PostPart(part *data.Part) errorCode {
	client := &http.Client{}
	data, err := json.Marshal(part)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest
	}

	req, err := http.NewRequest("POST", managSysMoel.PartsEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status from insertion received: %s", http.StatusText(resp.StatusCode))
		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (managSysModel *ManagementSystemModel) GetAllOrders() (*[]data.Order, error) {
	resp, err := http.Get(managSysModel.OrdersEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ordersResp OrdersResponse

	err = json.Unmarshal(data, &ordersResp)
	if err != nil {
		return nil, err
	}

	return &ordersResp.Orders, nil
}

func (managSysModel *ManagementSystemModel) GetOrdersByClientId(clientId string) ([]data.Order, error) {
	resp, err := http.Get(managSysModel.OrdersEndpoint + "/search?clientid=" + clientId)
	if err != nil {
		return nil, fmt.Errorf("order with client ID %s not found", clientId)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status received: %s", resp.Status)
	}

	var ordersResponse OrdersResponse

	err = json.Unmarshal(data, &ordersResponse)
	if err != nil {
		return nil, err
	}

	return ordersResponse.Orders, nil
}

func (managSysMoel *ManagementSystemModel) PostOrder(order *data.Order) errorCode {
	client := &http.Client{}
	data, err := json.Marshal(order)
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest
	}

	req, err := http.NewRequest("POST", managSysMoel.OrdersEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Print(err)
		return http.StatusBadRequest
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return http.StatusInternalServerError
	}

	resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("unexpected status from insertion received: %s", http.StatusText(resp.StatusCode))
		return http.StatusInternalServerError
	}

	return http.StatusCreated
}

func (managSysMoel *ManagementSystemModel) RequestAuth(userAuth data.UserAuth) (bool, errorCode) {
	client := &http.Client{}

	data, err := json.Marshal(userAuth)
	if err != nil {
		log.Printf("[ERROR] - %s", err.Error())
		return false, http.StatusBadRequest
	}

	req, err := http.NewRequest("POST", managSysMoel.AuthEndpoint+"/", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("[ERROR] - %s", err.Error())
		return false, http.StatusInternalServerError
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] - %s", err.Error())
		return false, http.StatusInternalServerError
	}

	defer resp.Body.Close()

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] - %s", err.Error())
		return false, http.StatusInternalServerError
	}

	var authResponse AuthResponse
	err = json.Unmarshal(data, &authResponse)
	if err != nil {
		log.Printf("[%s - ERROR] %s", shared.GetCallerInfo(), err.Error())
		return false, http.StatusInternalServerError
	}

	log.Printf("[%s] Authorized: %v", shared.GetCallerInfo(), authResponse.Authorized)

	return authResponse.Authorized, http.StatusOK
}
