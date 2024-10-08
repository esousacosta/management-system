package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"

	"github.com/esousacosta/managementsystem/cmd/shared"
	"github.com/esousacosta/managementsystem/internal/appmodel"
)

type application struct {
	managSysModel appmodel.ManagementSystemModel
}

func main() {
	addr := flag.String("addr", "localhost:3000", "HTTP network address")
	partsEndpoint := flag.String("partsEndpoint", "https://localhost:4000/v1/parts", "Parts endpoint for accessing the management system web service")
	ordersEndpoint := flag.String("ordersEndpoint", "https://localhost:4000/v1/orders", "Orders endpoint for accessing the management system web service")
	authEndpoint := flag.String("authEndpoint", "https://localhost:4000/v1/auth", "Orders endpoint for accessing the management system web service")

	// cookieJar, err := cookiejar.New(nil)
	// if err != nil {
	// 	log.Fatalf("[%s] ERROR - %v", shared.GetCallerInfo(), err)
	// }

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: shared.GetCertPool()},
		},
		Jar: nil,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			log.Printf("[REDIRECT] received headers: %v", r.Header)
			r.Header.Add("Cookie", via[0].Header.Get("Cookie"))
			return nil
		},
	}

	app := application{managSysModel: appmodel.NewManagementSystemModel(*ordersEndpoint, *partsEndpoint, *authEndpoint, client)}

	cert, err := tls.LoadX509KeyPair("./cert/domain.crt", "./cert/private.key")
	if err != nil {
		log.Fatalf("[%s - ERROR] %s", shared.GetCallerInfo(), err)
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}, RootCAs: shared.GetCertPool(), InsecureSkipVerify: false, ServerName: "localhost"}

	srv := http.Server{
		Addr:      *addr,
		TLSConfig: tlsCfg,
		Handler:   app.route(),
	}

	log.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("cert/domain.crt", "cert/private.key")
	log.Fatal(err)
}
