package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	invoiceModel "github.com/semka95/balance-service/invoice/repository"
	transferModel "github.com/semka95/balance-service/transfer/repository"
	userModel "github.com/semka95/balance-service/user/repository"
)

// API represents rest api
type API struct {
	userStore     userModel.Querier
	transferStore transferModel.Querier
	invoiceStore  invoiceModel.Querier
	db            *sql.DB
}

// NewRouter creates api router
func (a *API) NewRouter(userStore userModel.Querier, tranferStore transferModel.Querier, invoiceStore invoiceModel.Querier, db *sql.DB) chi.Router {
	a.userStore = userStore
	a.transferStore = tranferStore
	a.invoiceStore = invoiceStore
	a.db = db

	r := chi.NewRouter()
	r.Route("/api/v1/user", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getBalance)
		rapi.Put("/deposit", a.depositMoney)
		rapi.Put("/withdraw", a.withdrawMoney)
		rapi.Post("/", a.createUser)
		rapi.Put("/transfer", a.transfer)
	})
	r.Route("/api/v1/transfer", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getTransfer)
		rapi.Get("/{user_id}/inbound", a.getInboundTransfers)
		rapi.Get("/{user_id}/outbound", a.getOutboundTransfers)
		rapi.Get("/{from_uid}/to/{to_uid}", a.getTransfersBetweenUsers)
	})
	r.Route("/api/v1/invoice", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getInvoice)
		rapi.Get("/user/{id}", a.getInvoiceByUserID)
		rapi.Post("/", a.createInvoice)
		rapi.Put("/{id}/accept", a.acceptInvoice)
	})

	return r
}
