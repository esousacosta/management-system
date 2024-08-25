package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/internal/appmodel"
)

type application struct {
	managSysModel appmodel.ManagementSystemModel
}

func main() {
	addr := flag.String("addr", ":80", "HTTP network address")
	endpoint := flag.String("endpoint", "http://localhost:4000/v1/parts", "Endpoint for accessing the management system web service")

	app := application{managSysModel: appmodel.NewManagementSystemModel(*endpoint)}

	srv := http.Server{
		Addr:    *addr,
		Handler: app.route(),
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
