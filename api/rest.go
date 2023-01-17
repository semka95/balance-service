package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"github.com/semka95/balance-service/domain"
	invoiceModel "github.com/semka95/balance-service/invoice/repository"
	transferModel "github.com/semka95/balance-service/transfer/repository"
	userModel "github.com/semka95/balance-service/user/repository"
)

// API represents rest api
type API struct {
	userStore     userModel.Querier
	userUcase     domain.UserUsecase
	transferStore transferModel.Querier
	transferUcase domain.TransferUsecase
	invoiceStore  invoiceModel.Querier
	db            *sql.DB
}

// NewRouter creates api router
func (a *API) NewRouter(userStore userModel.Querier, userUcase domain.UserUsecase, tranferStore transferModel.Querier, transferUcase domain.TransferUsecase, invoiceStore invoiceModel.Querier, db *sql.DB) chi.Router {
	a.userStore = userStore
	a.userUcase = userUcase
	a.transferStore = tranferStore
	a.transferUcase = transferUcase
	a.invoiceStore = invoiceStore
	a.db = db

	r := chi.NewRouter()
	r.Route("/api/v1/user", func(rapi chi.Router) {
		rapi.Get("/{id}", a.getBalance)
		rapi.Put("/{id}/deposit", a.depositMoney)
		rapi.Put("/{id}/withdraw", a.withdrawMoney)
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
