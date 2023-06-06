package main

import (
	"fmt"
	"github.com/angelmotta/flow-currency-api/internal/exchangestore"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
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
	// Create a new CurrencyServer
	currServer := NewCurrencyServer()

	// Retrieve data from DB
	val, err := currServer.rdb.GetExchange("usd_pen_")
	if err != nil {
		log.Fatalf("error getting exchange rate: %v", err)
	} else if val == "" {
		log.Println("currency exchange does not exist")
	}
	log.Println("currency exchange: ", val)

	// Create a new router
	r := chi.NewRouter()

	r.Get("/exchanges", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("currency"))
		if err != nil {
			return
		}
	})

	r.Get("/exchanges/{idSrcCurrency}", func(w http.ResponseWriter, r *http.Request) {
		idSrcCurr := chi.URLParam(r, "idSrcCurrency")
		log.Println("idSrcCurrency:", idSrcCurr)
		_, err := fmt.Fprintf(w, "idSrcCurrency: %s\n", idSrcCurr)
		if err != nil {
			return
		}
	})

	r.Get("/exchanges/{idSrcCurrency}/{idDstCurrency}", func(w http.ResponseWriter, r *http.Request) {
		idSrcCurr := chi.URLParam(r, "idSrcCurrency")
		idDstCurr := chi.URLParam(r, "idDstCurrency")
		log.Println("idSrcCurrency:", idSrcCurr)
		log.Println("idDstCurrency:", idDstCurr)
		_, err := fmt.Fprintf(w, "idSrcCurrency: %s, idDstCurrency: %s\n", idSrcCurr, idDstCurr)
		if err != nil {
			return
		}
	})

	// Start HTTP server
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Panicf("error starting server: %v", err)
	}
}
