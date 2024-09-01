package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/esousacosta/managementsystem/cmd/shared"
)

type UserAuth struct {
	ID        int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
}

type ReadUserAuth struct {
	ID        int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
	Email     *string   `json:"email"`
	Password  *string   `json:"password"`
}

type UserAuthModel struct {
	db *sql.DB
}

func (um *UserAuthModel) GetUserAuth(email string) (*UserAuth, error) {
	query := `SELECT *
				FROM user_auth
				WHERE email = $1`

	var user UserAuth

	err := um.db.QueryRow(query, email).Scan(&user.Email, &user.Password, &user.CreatedAt, &user.ID)
	if err != nil {
		log.Printf("[%s] get user_auth info query error --> %s", shared.GetCallerInfo(), err.Error())
		return nil, fmt.Errorf("error fetching user information from the database")
	}

	return &user, nil
}

func (um *UserAuthModel) InsertUser(userAuth *ReadUserAuth) error {
	query := `INSERT INTO user_auth (email, password)
				VALUES ($1, $2)
				RETURNING id, created_at`

	args := []any{*userAuth.Email, *userAuth.Password}

	err := um.db.QueryRow(query, args...).Scan(&userAuth.ID, &userAuth.CreatedAt)
	if err != nil {
		log.Print("user auth insert error --> " + err.Error())
		return fmt.Errorf("error creating user auth information in the database")
	}

	return nil
}
