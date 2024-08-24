package data

import (
	"database/sql"
	"time"
)

type part struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	Stock     int64     `json:"stock"`
	Reference string    `json:"reference"`
	BarCode   string    `json:"barcode"`
}

type PartModel struct {
	db *sql.DB
}
type PartsReponse struct {
	Parts *[]part `json:"parts"`
}

func (partModel *PartModel) GetAll() ([]*part, error) {
	query := `SELECT *
	FROM parts
	ORDER BY id;`

	rows, err := partModel.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var parts []*part

	for rows.Next() {
		var p part
		err = rows.Scan(&p.Id, &p.CreatedAt, &p.Name, &p.Price, &p.Stock, &p.Reference, &p.BarCode)
		if err != nil {
			return nil, err
		}
		parts = append(parts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return parts, err
}
