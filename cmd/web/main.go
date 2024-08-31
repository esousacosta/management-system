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
	addr := flag.String("addr", "localhost:3000", "HTTP network address")
	partsEndpoint := flag.String("partsEndpoint", "http://localhost:4000/v1/parts", "Parts endpoint for accessing the management system web service")
	ordersEndpoint := flag.String("ordersEndpoint", "http://localhost:4000/v1/orders", "Orders endpoint for accessing the management system web service")
	authEndpoint := flag.String("authEndpoint", "http://localhost:4000/v1/auth", "Orders endpoint for accessing the management system web service")

	app := application{managSysModel: appmodel.NewManagementSystemModel(*ordersEndpoint, *partsEndpoint, *authEndpoint)}

	srv := http.Server{
		Addr:    *addr,
		Handler: app.route(),
	}

	log.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	log.Fatal(err)
}
