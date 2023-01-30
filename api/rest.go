package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"github.com/semka95/balance-service/domain"
)

// API represents rest api
type API struct {
	userUcase     domain.UserUsecase
	transferUcase domain.TransferUsecase
	invoiceUcase  domain.InvoiceUsecase
	db            *sql.DB
}

// NewRouter creates api router
func (a *API) NewRouter(userUcase domain.UserUsecase, transferUcase domain.TransferUsecase, invoiceUcase domain.InvoiceUsecase, db *sql.DB) chi.Router {
	a.userUcase = userUcase
	a.transferUcase = transferUcase
	a.invoiceUcase = invoiceUcase
	a.db = db

	r := chi.NewRouter()
	r.Route("/api/v1/user", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getBalance)
		rapi.Patch("/{id}/deposit", a.depositMoney)
		rapi.Patch("/{id}/withdraw", a.withdrawMoney)
		rapi.Post("/", a.createUser)
	})
	r.Route("/api/v1/transfer", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getTransfer)
		rapi.Get("/{user_id}/inbound", a.getInboundTransfers)
		rapi.Get("/{user_id}/outbound", a.getOutboundTransfers)
		rapi.Get("/{from_uid}/to/{to_uid}", a.getTransfersBetweenUsers)
		rapi.Post("/", a.createTransfer)
	})
	r.Route("/api/v1/invoice", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getInvoice)
		rapi.Get("/user/{id}", a.getInvoiceByUserID)
		rapi.Post("/", a.createInvoice)
		rapi.Put("/{id}/accept", a.acceptInvoice)
	})

	return r
}
