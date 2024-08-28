package appmodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/internal/data"
)

type errorCode int

type ManagementSystemModel struct {
	PartsEndpoint  string
	OrdersEndpoint string
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

func NewManagementSystemModel(ordersEndpoint string, partsEndpoint string) ManagementSystemModel {
	return ManagementSystemModel{
		PartsEndpoint:  partsEndpoint,
		OrdersEndpoint: ordersEndpoint,
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
