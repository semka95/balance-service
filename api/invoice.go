package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	invoiceModel "github.com/semka95/balance-service/invoice/repository"
)

// GET /invoice/{id} - returns invoice by id
func (a *API) getInvoice(w http.ResponseWriter, r *http.Request) {
	invID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid invoice id")
		return
	}

	invoice, err := a.invoiceStore.GetInvoiceByID(r.Context(), int64(invID))
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, fmt.Sprintf("invoice with %d id not found", invID))
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get invoice")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, invoice)
}

// GET /invoice/user/{id} - returns invoices by user id
func (a *API) getInvoiceByUserID(w http.ResponseWriter, r *http.Request) {
	uID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 10
	}
	cursor, err := strconv.Atoi(r.URL.Query().Get("cursor"))
	if err != nil {
		cursor = 0
	}

	params := invoiceModel.GetInvoicesByUserIDParams{
		UserID: int64(uID),
		ID:     int64(cursor),
		Limit:  int32(limit),
	}

	invoices, err := a.invoiceStore.GetInvoicesByUserID(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, fmt.Sprintf("invoice with %d id not found", uID))
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get invoice")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, invoices)
}

// POST /invoice - create invoice
func (a *API) createInvoice(w http.ResponseWriter, r *http.Request) {
	params := invoiceModel.CreateInvoiceParams{}

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to invoice")
		return
	}

	invoice, err := a.invoiceStore.CreateInvoice(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't invoice user")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, invoice)
}

// PUT /invoice/{id}/accept - accept invoice
func (a *API) acceptInvoice(w http.ResponseWriter, r *http.Request) {
	invID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid invoice id")
		return
	}

	tx, err := a.db.BeginTx(r.Context(), nil)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't start transaction")
		return
	}
	defer tx.Rollback()

	params := invoiceModel.UpdateStatusParams{
		ID:            int64(invID),
		PaymentStatus: invoiceModel.ValidStatusAccepted,
	}
	_, err = a.invoiceStore.GetInvoiceByID(r.Context(), params.ID)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "invoice not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get invoice")
		return
	}

	rows, err := a.invoiceStore.UpdateStatus(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update invoice")
		return
	}

	if err := tx.Commit(); err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't commit transaction")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, rows)
}
