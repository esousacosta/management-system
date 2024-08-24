package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/internal/data"
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

func readJson(w http.ResponseWriter, r *http.Request, logger *log.Logger) (*data.Part, error) {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var input data.Part
	if err := dec.Decode(&input); err != nil {
		logger.Printf("decoding error --> %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, err
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, errors.New("body must contain only one JSON object")
	}

	return &input, nil
}
