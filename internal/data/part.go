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
}

type PartModel struct {
	db *sql.DB
}

func (partModel *PartModel) GetAll() {

}
