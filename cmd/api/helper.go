package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type envelope map[string]any

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

func writeJson(w http.ResponseWriter, statusCode int, data envelope, headers http.Header) error {
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return err
	}

	json = append(json, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)

	return nil
}
