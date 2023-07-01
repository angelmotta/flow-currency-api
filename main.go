package main

import (
	"github.com/angelmotta/flow-currency-api/app"
	"log"
	"net/http"
)

func main() {
	// Create a new CurrencyServer
	currServer := app.NewCurrencyServer()
	currServer.Routes()

	// Start HTTP server
	log.Fatal(http.ListenAndServe(":8080", currServer.Router))
}
