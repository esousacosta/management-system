package data

import (
	"database/sql"
	"time"
)

type OrderModel struct {
	db *sql.DB
}

type Order struct {
	ID        int       `json:"-"`
	ClientId  string    `json:"clientid"`
	CreatedAt time.Time `json:"-"`
	Services  []string  `json:"services"`
	PartsIds  []int     `json:"partsids"`
	Comment   string    `json:"comment"`
	Total     float32   `json:"total"`
}

// func (orderModel *OrderModel) GetAll() []*Order {
// }
