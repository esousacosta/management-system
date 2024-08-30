package data

import (
	"database/sql"
	"fmt"
	"time"
)

type UserAuth struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type UserAuthModel struct {
	db *sql.DB
}

func (um *UserAuthModel) GetUser(email string) (*UserAuth, error) {
	query := `SELECT *
				FROM users
				WHERE email = $1`

	var user UserAuth

	err := um.db.QueryRow(query, email).Scan(&user.ID, &user.CreatedAt, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("error fetching user information from the database")
	}

	return &user, nil
}
