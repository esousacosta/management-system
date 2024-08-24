package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func getDbConnection(dsn *string, logger *log.Logger) *sql.DB {
	if dsn == nil {
		logger.Println("invalid DSN received. Connection with the db can't be stablished")
		return nil
	}
	var db *sql.DB

	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	logger.Println("Connection with the database successfully stablished.")

	return db
}
