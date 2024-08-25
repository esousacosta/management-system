package data

import (
	"database/sql"
	"time"
)

type Part struct {
	Id        int       `json:"-"`
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
	Parts []*Part `json:"parts"`
}

func (partModel *PartModel) GetAll() ([]*Part, error) {
	query := `SELECT *
	FROM parts
	ORDER BY id;`

	rows, err := partModel.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var parts []*Part

	for rows.Next() {
		var p Part
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

func (partModel *PartModel) GetById(id int64) (*Part, error) {
	query := `SELECT * FROM parts
				WHERE id = $1`

	row := partModel.db.QueryRow(query, id)

	var part Part

	err := row.Scan(&part.Id, &part.CreatedAt, &part.Name, &part.Price, &part.Stock, &part.Reference, &part.BarCode)
	if err != nil {
		return nil, err
	}

	return &part, nil
}

func (partModel *PartModel) GetByRef(ref string) (*Part, error) {
	query := `SELECT * FROM parts
				WHERE reference = $1`

	row := partModel.db.QueryRow(query, ref)

	var part Part

	err := row.Scan(&part.Id, &part.CreatedAt, &part.Name, &part.Price, &part.Stock, &part.Reference, &part.BarCode)
	if err != nil {
		return nil, err
	}

	return &part, nil
}

func (partModel *PartModel) Insert(part *Part) error {
	query := `INSERT INTO parts (name, price, stock, reference, barcode)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, created_at`

	// interface{} === any
	args := []interface{}{part.Name, part.Price, part.Stock, part.Reference, part.BarCode}

	return partModel.db.QueryRow(query, args...).Scan(&part.Id, &part.CreatedAt)
}
