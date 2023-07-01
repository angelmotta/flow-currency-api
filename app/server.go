package app

import (
	"encoding/json"
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

type exchangeResponse struct {
	SrcCurrency string `json:"srcCurrency"`
	DstCurrency string `json:"dstCurrency"`
	Rate        string `json:"rate"`
}

type errorResponse struct {
	Details string `json:"details"`
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

	// Get query params (Eg. amount=100)
	receivedAmount := r.URL.Query().Get("amount")
	if receivedAmount != "" {
		log.Println("receivedAmount:", receivedAmount)
		// TODO: Prepare response -> convert string to float64 and multiply by the exchange rate
	}

	// Get Exchange Rate from DB
	pairCurrency := idSrcCurr + "_" + idDstCurr
	exchangeRate, err := cs.Rdb.GetExchange(pairCurrency)
	if err != nil {
		http.Error(w, "Error DB service", http.StatusInternalServerError)
		return
	}
	var response interface{}
	var statusCode int
	// Verify if exchangeRate is empty
	if exchangeRate != "" {
		response = exchangeResponse{
			SrcCurrency: idSrcCurr,
			DstCurrency: idDstCurr,
			Rate:        exchangeRate,
		}
		statusCode = http.StatusOK
	} else {
		response = errorResponse{
			Details: fmt.Sprintf("Exchange rate not found: '%v' to '%v'", idSrcCurr, idDstCurr),
		}
		statusCode = http.StatusNotFound
	}

	renderJsonUtil(w, response, statusCode)
}

func renderJsonUtil(w http.ResponseWriter, payload interface{}, statusCode int) {
	responseJson, err := json.Marshal(payload)
	if err != nil {
		// Marshal error (internal server error)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	// Write Json response
	_, err = w.Write(responseJson)
	if err != nil {
		log.Printf("Error sending response: %v", err.Error())
		return
	}
}
