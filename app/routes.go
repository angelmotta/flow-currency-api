package app

func (cs *CurrencyServer) Routes() {
	cs.Router.Get("/api/exchanges", cs.GetAllExchangesHandler)
	cs.Router.Get("/api/exchanges/{idDstCurrency}", cs.GetExchangesHandler)
	cs.Router.Get("/api/exchanges/{idDstCurrency}/{idSrcCurrency}", cs.GetExchangeHandler)
}
