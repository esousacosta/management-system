package data

import (
	"database/sql"
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
	CreatedAt time.Time `json:"created_at"`
	Services  []string  `json:"services"`
	PartsRefs []string  `json:"parts_refs"`
	Comment   string    `json:"comment"`
	Total     float32   `json:"total"`
}

type ReadOrder struct {
	ID        int       `json:"-"`
	ClientId  *string   `json:"client_id"`
	CreatedAt time.Time `json:"created_at"`
	Services  *[]string `json:"services"`
	PartsRefs *[]string `json:"parts_refs"`
	Comment   *string   `json:"comment"`
	Total     *float32  `json:"total"`
}

func (om *OrderModel) GetAll() ([]*Order, error) {
	query := `SELECT *
			FROM orders
			ORDER BY id`

	rows, err := om.db.Query(query)
	if err != nil {
		log.Print("query error: ")
		return nil, err
	}

	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.ClientId, &order.CreatedAt, pq.Array(&order.Services), pq.Array(&order.PartsRefs), &order.Comment, &order.Total); err != nil {
			log.Print("scan error: ")
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (om *OrderModel) Get(orderId int64) (*Order, error) {
	query := `SELECT *
			FROM orders
			WHERE id = $1`

	var order Order

	err := om.db.QueryRow(query, orderId).Scan(&order.ID, &order.ClientId, &order.CreatedAt, pq.Array(&order.Services), pq.Array(&order.PartsRefs), &order.Comment, &order.Total)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil, err
	}

	return &order, nil
}

func (om *OrderModel) GetByClientId(clientId int64) ([]*Order, error) {
	query := `SELECT *
			FROM orders
			WHERE client_id = $1`

	rows, err := om.db.Query(query, clientId)
	if err != nil {
		log.Printf("query error: %v", err)
		return nil, err
	}

	defer rows.Close()

	var orders []*Order

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.ClientId, &order.CreatedAt, pq.Array(&order.Services), pq.Array(&order.PartsRefs), &order.Comment, &order.Total)
		if err != nil {
			log.Print("scan error: ")
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (om *OrderModel) Insert(orderToInsert *ReadOrder) error {
	query := `INSERT INTO orders (client_id, services, parts_refs, comment, total)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

	args := []any{orderToInsert.ClientId, pq.Array(orderToInsert.Services), pq.Array(orderToInsert.PartsRefs), orderToInsert.Comment, orderToInsert.Total}

	err := om.db.QueryRow(query, args...).Scan(&orderToInsert.ID, &orderToInsert.CreatedAt)
	if err != nil {
		log.Printf("insert query error --> %v", err.Error())
		return err
	}

	return nil
}

func (om *OrderModel) Update(order *Order) error {
	query := `UPDATE orders SET
				client_id = $1,
				services = $2,
				parts_refs = $3,
				comment = $4,
				total = $5
				WHERE id = $6
				RETURNING client_id`

	args := []any{order.ClientId, pq.Array(&order.Services), pq.Array(order.PartsRefs), order.Comment, order.Total, order.ID}
	result, err := om.db.Exec(query, args...)
	if err != nil {
		log.Printf("error updating entry in orders db: %v", err)
		return err
	}

	if nbRowsAffected, err := result.RowsAffected(); nbRowsAffected == 0 || err != nil {
		log.Printf("no orders updated for id %d", order.ID)
		return err
	}

	return nil
}

func (om *OrderModel) Delete(orderId int64) error {
	query := `DELETE FROM orders
				WHERE id = $1`

	result, err := om.db.Exec(query, orderId)
	if err != nil {
		log.Printf("error updating entry in orders db: %v", err)
		return err
	}

	if nbRowsAffected, err := result.RowsAffected(); nbRowsAffected == 0 || err != nil {
		log.Printf("no orders delete for id %d", orderId)
		return err
	}

	return nil
}
