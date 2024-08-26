package data

import (
	"database/sql"
	"fmt"
	"time"
)

type OrderModel struct {
	db *sql.DB
}

type Order struct {
	ID        int       `json:"-"`
	ClientId  string    `json:"client_id"`
	CreatedAt time.Time `json:"-"`
	Services  []string  `json:"services"`
	PartsIds  []int     `json:"parts_ids"`
	Comment   string    `json:"comment"`
	Total     float32   `json:"total"`
}

func (om *OrderModel) GetAll() ([]*Order, error) {
	query := `SELECT *
			FROM orders
			ORDER BY id`

	rows, err := om.db.Query(query)
	if err != nil {
		fmt.Print("query error: ")
		return nil, err
	}

	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientId, &order.PartsIds, &order.PartsIds, &order.PartsIds, &order.Comment, &order.Total); err != nil {
			fmt.Print("scan error: ")
			return nil, err
		}
		fmt.Printf("%v", order)
		orders = append(orders, &order)
	}

	return orders, nil
}

// func (orderModel *OrderModel) GetAll() []*Order {
// }
