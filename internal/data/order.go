package data

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type OrderModel struct {
	db *sql.DB
}

type Order struct {
	ID        int       `json:"-"`
	ClientId  string    `json:"client_id"`
	CreatedAt time.Time `json:"-"`
	Services  []string  `json:"services"`
	PartsIds  []int64   `json:"parts_ids"`
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
		// the drive.Value, the thing Scan uses to read values,
		// doesn't parse int slices - hence the need for the hack below.
		var partsIdsArr pq.Int64Array
		if err := rows.Scan(&order.ClientId, &order.CreatedAt, pq.Array(order.Services), &partsIdsArr, &order.Comment, &order.Total, &order.ID); err != nil {
			fmt.Print("scan error: ")
			return nil, err
		}
		order.PartsIds = []int64(partsIdsArr)
		fmt.Printf("%+v", order)
		orders = append(orders, &order)
	}

	return orders, nil
}

// func (orderModel *OrderModel) GetAll() []*Order {
// }
