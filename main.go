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
	log.Println("Running flowApp-currencyServer at port 8081")
	log.Fatal(http.ListenAndServe(":8081", currServer.Router))
}
