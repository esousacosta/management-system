package data

import "database/sql"

type Models struct {
	Parts PartModel
}

func NewModel(db *sql.DB) *Models {
	return &Models{
		Parts: PartModel{
			db: db,
		},
	}
}
