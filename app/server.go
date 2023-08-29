package app

import (
	"encoding/json"
	"fmt"
	"github.com/angelmotta/flow-currency-api/internal/exchangestore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"log"
	"math"
	"net/http"
	"strconv"
)

// CurrencyServer is the main struct for the Currency API service
type CurrencyServer struct {
	Rdb    *exchangestore.ExchangeStore
	Router *chi.Mux
}

type exchangeResponse struct {
	SrcCurrency   string  `json:"srcCurrency"`
	DstCurrency   string  `json:"dstCurrency"`
	Rate          string  `json:"rate"`
	MoneyReceived string  `json:"moneyReceived,omitempty"`
	MoneySent     float64 `json:"moneySent,omitempty"`
}

type errorResponse struct {
	Details string `json:"details"`
}

// NewCurrencyServer creates a new CurrencyServer
func NewCurrencyServer() *CurrencyServer {
	// Create a new DB connection
	redisAddr := "localhost:6379"
	dbConn, err := exchangestore.New(redisAddr)
	if err != nil {
		log.Fatalf("error creating redis client: %v", err)
	}
	// Create a new http Router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://mysideproject.com", "http://localhost:5173"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		//AllowCredentials: false,
		//MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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

// GetExchangesHandler retrieves all available currency exchanges for the given currency
func (cs *CurrencyServer) GetExchangesHandler(w http.ResponseWriter, r *http.Request) {
	idSrcCurr := chi.URLParam(r, "idSrcCurrency")
	log.Println("idSrcCurrency:", idSrcCurr)
	_, err := fmt.Fprintf(w, "idSrcCurrency: %s\n", idSrcCurr)
	if err != nil {
		return
	}
}

// GetExchangeHandler retrieves a specific currency exchange
func (cs *CurrencyServer) GetExchangeHandler(w http.ResponseWriter, r *http.Request) {
	idSrcCurr := chi.URLParam(r, "idSrcCurrency")
	idDstCurr := chi.URLParam(r, "idDstCurrency")
	log.Println("idSrcCurrency:", idSrcCurr)
	log.Println("idDstCurrency:", idDstCurr)

	// Get Exchange Rate from DB
	pairCurrency := idDstCurr + "_" + idSrcCurr // usd_pen (casa vende USD)
	exchangeRate, err := cs.Rdb.GetExchange(pairCurrency)
	if err != nil {
		http.Error(w, "Error DB service", http.StatusInternalServerError)
		return
	}
	// Check if 'amount' exist in query param
	receivedAmount := r.URL.Query().Get("amount")
	var moneyConverted float64
	if receivedAmount != "" {
		log.Println("receivedAmount:", receivedAmount)
		// convert string receivedAmount to float64
		receivedMoney, err := strconv.ParseFloat(receivedAmount, 64)
		if err != nil {
			responseError := errorResponse{
				Details: fmt.Sprintf("Exchange rate not found: '%v' to '%v'", idSrcCurr, idDstCurr),
			}
			renderJsonUtil(w, responseError, http.StatusBadRequest)
			log.Println("test error casting")
			return
		}
		rate, err := strconv.ParseFloat(exchangeRate, 64)
		moneyConverted = receivedMoney / rate
		// Round
		moneyConverted = math.Round(moneyConverted*100) / 100
	}

	var response interface{}
	var statusCode int
	// Verify if exchangeRate is not empty
	if exchangeRate != "" {
		response = exchangeResponse{
			SrcCurrency:   idSrcCurr,
			DstCurrency:   idDstCurr,
			Rate:          exchangeRate,
			MoneyReceived: receivedAmount,
			MoneySent:     moneyConverted,
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
