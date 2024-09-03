package data

import "database/sql"

type Models struct {
	Clients   ClientModel
	Parts     PartModel
	Orders    OrderModel
	UsersAuth UserAuthModel
}

func NewModel(db *sql.DB) *Models {
	return &Models{
		Clients: ClientModel{
			db: db,
		},
		Parts: PartModel{
			db: db,
		},
		Orders: OrderModel{
			db: db,
		},
		UsersAuth: UserAuthModel{
			db: db,
		},
	}
}
