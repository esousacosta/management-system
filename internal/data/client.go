package data

import (
	"database/sql"
	"time"
)

type Client struct {
	Id        int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Reference string    `json:"reference"`
	UserId    int       `json:"-"`
}

type ReadClient struct {
	Id        int       `json:"-"`
	CreatedAt time.Time `json:"-"`
	Name      *string   `json:"name"`
	LastName  *string   `json:"lastname"`
	Email     *string   `json:"email"`
	Phone     *string   `json:"phone"`
	Reference *string   `json:"reference"`
	UserId    *int      `json:"-"`
}

type ClientModel struct {
	db *sql.DB
}

func (cm *ClientModel) GetAll() ([]*Client, error) {
	query := `SELECT *
				FROM clients`

	rows, err := cm.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var clients []*Client

	for rows.Next() {
		var client Client
		err := rows.Scan(&client.Id, &client.CreatedAt, &client.Name, &client.LastName, &client.Email, &client.Phone, &client.Reference, &client.UserId)
		if err != nil {
			return nil, err
		}
		clients = append(clients, &client)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (cm *ClientModel) GetClientById(clientId int) (*Client, error) {
	query := `SELECT *
				FROM clients
				WHERE id = $1`

	row := cm.db.QueryRow(query, clientId)

	var client Client

	err := row.Scan(&client.Id, &client.CreatedAt, &client.Name, &client.LastName, &client.Email, &client.Phone, &client.Reference, &client.UserId)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (cm *ClientModel) Insert(client *ReadClient, userId int) error {
	query := `INSERT INTO clients (name, lastname, email, phone, reference, user_id)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id, created_at`

	args := []interface{}{client.Name, client.LastName, client.Email, client.Phone, client.Reference, userId}

	return cm.db.QueryRow(query, args...).Scan(&client.Id, &client.CreatedAt)
}
