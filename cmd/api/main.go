package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/esousacosta/managementsystem/internal/data"
)

type config struct {
	port int
	dsn  string
	env  string
}

type application struct {
	config config
	logger *log.Logger
	model  *data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.env, "env", "dev", "Server Environment")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("MANAGEMENT_SYSTEM_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db := getDbConnection(&cfg.dsn, logger)

	app := &application{
		config: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		model:  data.NewModel(db),
	}

	addr := fmt.Sprintf("localhost:%d", app.config.port)

	cert, err := tls.LoadX509KeyPair("./cert/domain.crt", "./cert/private.key")
	if err != nil {
		log.Fatalf("[%s - ERROR] %s", shared.GetCallerInfo(), err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{cert}}
	// InsecureSkipVerify: true}

	srv := http.Server{
		Addr:         addr,
		TLSConfig:    tlsCfg,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("[%s] Server starting on port %d", app.config.env, app.config.port)
	err = srv.ListenAndServeTLS("", "")
	app.logger.Fatal(err)
}
