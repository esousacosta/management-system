package appmodel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/esousacosta/managementsystem/internal/data"
)

type ManagementSystemModel struct {
	Endpoint string
}

type PartsResponse struct {
	Parts []data.Part `json:"parts"`
}

type PartResponse struct {
	Parts data.Part `json:"part"`
}

func NewManagementSystemModel(endpoint string) ManagementSystemModel {
	return ManagementSystemModel{
		Endpoint: endpoint,
	}
}

func (managSysModel *ManagementSystemModel) GetAll() (*[]data.Part, error) {
	resp, err := http.Get(managSysModel.Endpoint)
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
	resp, err := http.Get(managSysModel.Endpoint + "/" + partRef)
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
