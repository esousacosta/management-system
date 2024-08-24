package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	dsn  string
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.env, "env", "dev", "Server Environment")
	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("MANAGEMENT_SYSTEM_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	addr := fmt.Sprintf(":%d", app.config.port)

	srv := http.Server{
		Addr:         addr,
		Handler:      app.route(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("[%s] Server starting on port %d", app.config.env, app.config.port)
	err := srv.ListenAndServe()
	app.logger.Fatal(err)
}
