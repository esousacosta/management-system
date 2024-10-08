package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Part struct {
	Id        int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Price     float32   `json:"price"`
	Stock     int64     `json:"stock"`
	Reference string    `json:"reference"`
	Barcode   string    `json:"barcode"`
	UserId    int       `json:"-"`
}

type ReadPart struct {
	Id        int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	Name      *string   `json:"name"`
	Price     *float32  `json:"price"`
	Stock     *int64    `json:"stock"`
	Reference *string   `json:"reference"`
	Barcode   *string   `json:"barcode"`
	UserId    *int      `json:"user_id"`
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
		err = rows.Scan(&p.Id, &p.CreatedAt, &p.Name, &p.Price, &p.Stock, &p.Reference, &p.Barcode, &p.UserId)
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

	err := row.Scan(&part.Id, &part.CreatedAt, &part.Name, &part.Price, &part.Stock, &part.Reference, &part.Barcode, &part.UserId)
	if err != nil {
		return nil, err
	}

	return &part, nil
}

func (partModel *PartModel) GetByRef(ref string) (*Part, error) {
	query := `SELECT * FROM parts
				WHERE LOWER(reference) = LOWER($1)`

	row := partModel.db.QueryRow(query, ref)

	var part Part

	err := row.Scan(&part.Id, &part.CreatedAt, &part.Name, &part.Price, &part.Stock, &part.Reference, &part.Barcode, &part.UserId)
	if err != nil {
		return nil, errors.New("part with requested reference not found")
	}

	return &part, nil
}

func (partModel *PartModel) Insert(part *ReadPart) error {
	query := `INSERT INTO parts (name, price, stock, reference, barcode, user_id)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id, created_at`

	// interface{} === any
	args := []interface{}{part.Name, part.Price, part.Stock, part.Reference, part.Barcode, part.UserId}

	return partModel.db.QueryRow(query, args...).Scan(&part.Id, &part.CreatedAt)
}

func (partModel *PartModel) Update(part *Part) error {
	query := `UPDATE parts SET
				name = $1,
				price = $2,
				stock = $3,
				reference = $4,
				barcode = $5,
				user_id = $6
				WHERE reference = $7
				RETURNING name`

	args := []any{part.Name, part.Price, part.Stock, part.Reference, part.Barcode, part.UserId, part.Reference}
	_, err := partModel.db.Exec(query, args...)
	return err
}

func (partModel *PartModel) Delete(ref string) error {
	if ref == "" {
		return fmt.Errorf("no part found with ref %s", ref)
	}

	query := `DELETE FROM parts
				WHERE reference = $1`

	results, err := partModel.db.Exec(query, ref)
	if err != nil {
		return err
	}

	deletedRows, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if deletedRows == 0 {
		return fmt.Errorf("no part found with ref %s", ref)
	}

	return nil
}
