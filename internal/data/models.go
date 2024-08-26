package data

import "database/sql"

type Models struct {
	Parts  PartModel
	Orders OrderModel
}

func NewModel(db *sql.DB) *Models {
	return &Models{
		Parts: PartModel{
			db: db,
		},
		Orders: OrderModel{
			db: db,
		},
	}
}
