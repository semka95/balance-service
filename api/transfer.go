package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/lib/pq"

	transferModel "github.com/semka95/balance-service/transfer/repository"
)

// GET /transfer/{id} - returns transfer by id
func (a *API) getTransfer(w http.ResponseWriter, r *http.Request) {
	trID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid transfer id")
		return
	}

	transfer, err := a.transferUcase.GetTransfer(r.Context(), int64(trID))
	if errors.Is(err, sql.ErrNoRows) {
		render.Status(r, http.StatusNoContent)
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfer")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfer)
}

// GET /transfer/{user_id}/inbound?limit=5&cursor=0 - returns transfers that user received
func (a *API) getInboundTransfers(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
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

	params := transferModel.GetInboundTransfersParams{
		ToUserID: int64(userID),
		ID:       int64(cursor),
		Limit:    int32(limit),
	}

	transfers, err := a.transferUcase.GetInboundTransfers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		render.Status(r, http.StatusNoContent)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}

// GET /transfer/{user_id}/outbound?limit=5&cursor=0 - returns transfers that user sent
func (a *API) getOutboundTransfers(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
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

	params := transferModel.GetOutboundTransfersParams{
		FromUserID: int64(userID),
		ID:         int64(cursor),
		Limit:      int32(limit),
	}

	transfers, err := a.transferUcase.GetOutboundTransfers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		render.Status(r, http.StatusNoContent)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}

// GET /transfer/{from_uid}/to/{to_uid}?limit=5&cursor=0 - returns transfers between users
func (a *API) getTransfersBetweenUsers(w http.ResponseWriter, r *http.Request) {
	fromUID, err := strconv.Atoi(chi.URLParam(r, "from_uid"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid from user id")
		return
	}
	toUID, err := strconv.Atoi(chi.URLParam(r, "to_uid"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid to user id")
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

	params := transferModel.GetTransfersBetweenUsersParams{
		FromUserID: int64(fromUID),
		ToUserID:   int64(toUID),
		ID:         int64(cursor),
		Limit:      int32(limit),
	}

	transfers, err := a.transferUcase.GetTransfersBetweenUsers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		render.Status(r, http.StatusNoContent)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}

// POST /createTransfer - transfers money from one user to another
func (a *API) createTransfer(w http.ResponseWriter, r *http.Request) {
	params := transferModel.CreateTransferParams{}
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}

	transfer, err := a.transferUcase.CreateTransfer(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Constraint == "users_balance_check" {
		SendErrorJSON(w, r, http.StatusBadRequest, err, fmt.Sprintf("not enough money on '%d' balance", params.FromUserID))
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create transfer record")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfer)
}
