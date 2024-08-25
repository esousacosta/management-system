package appmodel

type ManagementSystemModel struct {
	Endpoint string
}

func NewManagementSystemModel(endpoint string) ManagementSystemModel {
	return ManagementSystemModel{
		Endpoint: endpoint,
	}
}
