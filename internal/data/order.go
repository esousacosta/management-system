package data

import (
	"database/sql"
	"fmt"
	"log"
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

type ReadOrder struct {
	ID        int       `json:"-"`
	ClientId  *string   `json:"client_id"`
	CreatedAt time.Time `json:"-"`
	Services  *[]string `json:"services"`
	PartsIds  *[]int64  `json:"parts_ids"`
	Comment   *string   `json:"comment"`
	Total     *float32  `json:"total"`
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
		if err := rows.Scan(&order.ID, &order.ClientId, &order.CreatedAt, pq.Array(order.Services), &partsIdsArr, &order.Comment, &order.Total); err != nil {
			fmt.Print("scan error: ")
			return nil, err
		}
		order.PartsIds = []int64(partsIdsArr)
		fmt.Printf("%+v", order)
		orders = append(orders, &order)
	}

	return orders, nil
}

func (om *OrderModel) Insert(orderToInsert *Order) error {
	query := `INSERT INTO orders (client_id, services, parts_ids, comment, total)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

	args := []any{orderToInsert.ClientId, pq.Array(orderToInsert.Services), pq.Int64Array(orderToInsert.PartsIds), orderToInsert.Comment, orderToInsert.Total}

	err := om.db.QueryRow(query, args...).Scan(&orderToInsert.ID, &orderToInsert.CreatedAt)
	if err != nil {
		log.Printf("insert query error --> %v", err.Error())
		return err
	}

	return nil
}

func (om *OrderModel) Get(orderId int64) (*Order, error) {
	query := `SELECT *
			FROM orders
			WHERE id = $1`

	var order Order

	var partsIdsArr pq.Int64Array
	err := om.db.QueryRow(query, orderId).Scan(&order.ID, &order.ClientId, &order.CreatedAt, pq.Array(&order.Services), &partsIdsArr, &order.Comment, &order.Total)
	if err != nil {
		fmt.Printf("query error: %v", err)
		return nil, err
	}

	order.PartsIds = []int64(partsIdsArr)

	return &order, nil
}
