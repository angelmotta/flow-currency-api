package app

import (
	"fmt"
	"github.com/angelmotta/flow-currency-api/internal/exchangestore"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// CurrencyServer is the main struct for the Currency API service
type CurrencyServer struct {
	Rdb    *exchangestore.ExchangeStore
	Router *chi.Mux
}

// NewCurrencyServer creates a new CurrencyServer
func NewCurrencyServer() *CurrencyServer {
	// Create a new DB connection
	redisAddr := "localhost:6380"
	dbConn, err := exchangestore.New(redisAddr)
	if err != nil {
		log.Fatalf("error creating redis client: %v", err)
	}
	// Create a new http Router
	r := chi.NewRouter()

	// Create and return a new CurrencyServer
	return &CurrencyServer{
		Rdb:    dbConn,
		Router: r,
	}
}

// GetAllExchangesHandler retrieves all currency exchanges
func (cs *CurrencyServer) GetAllExchangesHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Get All Currencies"))
	if err != nil {
		return
	}
}

// GetExchangesHandler retrieves all available currency exchanges for a given currency
func (cs *CurrencyServer) GetExchangesHandler(w http.ResponseWriter, r *http.Request) {
	idSrcCurr := chi.URLParam(r, "idSrcCurrency")
	log.Println("idSrcCurrency:", idSrcCurr)
	_, err := fmt.Fprintf(w, "idSrcCurrency: %s\n", idSrcCurr)
	if err != nil {
		return
	}
}

// GetExchangeHandler retrieves a specific currency exchange for a given currency
func (cs *CurrencyServer) GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	idSrcCurr := chi.URLParam(r, "idSrcCurrency")
	idDstCurr := chi.URLParam(r, "idDstCurrency")
	log.Println("idSrcCurrency:", idSrcCurr)
	log.Println("idDstCurrency:", idDstCurr)

	// Get query params (amount=100)
	receivedAmount := r.URL.Query().Get("amount")
	if receivedAmount != "" {
		log.Println("receivedAmount:", receivedAmount)
		// TODO: convert string to float64 and multiply by the exchange rate
	}

	_, err := fmt.Fprintf(w, "idSrcCurrency: %s, idDstCurrency: %s\n", idSrcCurr, idDstCurr)
	if err != nil {
		return
	}
}