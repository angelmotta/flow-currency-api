package app

func (cs *CurrencyServer) Routes() {
	cs.Router.Get("/exchanges", cs.GetAllExchangesHandler)

	cs.Router.Get("/exchanges/{idSrcCurrency}", cs.GetExchangesHandler)

	cs.Router.Get("/exchanges/{idSrcCurrency}/{idDstCurrency}", cs.GetExchangeHandler)
}
