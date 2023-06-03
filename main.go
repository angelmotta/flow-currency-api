package main

import (
	"github.com/angelmotta/flow-currency-api/internal/exchangestore"
	"log"
)

type CurrencyServer struct {
	rdb *exchangestore.ExchangeStore
}

func NewCurrencyServer() *CurrencyServer {
	redisAddr := "localhost:6380"
	dbConn, err := exchangestore.New(redisAddr)
	if err != nil {
		log.Fatalf("error creating redis client: %v", err)
	}
	return &CurrencyServer{rdb: dbConn}
}

func main() {
	log.Println("Hello world")

	// Create a new CurrencyServer
	cServer := NewCurrencyServer()
	val, err := cServer.rdb.GetExchange("usd_pen_")
	if err != nil {
		log.Fatalf("error getting exchange rate: %v", err)
	} else if val == "" {
		log.Println("currency exchange does not exist")
	}
	log.Println("currency exchange: ", val)
}
