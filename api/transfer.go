package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	transferModel "github.com/semka95/balance-service/transfer/repository"
)

// GET /transfer/{id} - returns transfer by id
func (a *API) getTransfer(w http.ResponseWriter, r *http.Request) {
	trID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid transfer id")
		return
	}

	transfer, err := a.transferStore.GetTransferByID(r.Context(), int64(trID))
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, fmt.Sprintf("transfer with %d id not found", trID))
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

	transfers, err := a.transferStore.GetInboundTransfers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("no inbound transfers was found for %d user id", userID), "no transfers found")
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

	transfers, err := a.transferStore.GetOutboundTransfers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("no outbound transfers was found for %d user id", userID), "no transfers found")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}

// GET /transfer/{from_uid}/to/{to_uid}?limit=5&cursor=0 - returns transfers between users
func (a *API) getTransfersBetweenUsers(w http.ResponseWriter, r *http.Request) {
	from_uid, err := strconv.Atoi(chi.URLParam(r, "from_uid"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid from user id")
		return
	}
	to_uid, err := strconv.Atoi(chi.URLParam(r, "to_uid"))
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
		FromUserID: int64(from_uid),
		ToUserID:   int64(to_uid),
		ID:         int64(cursor),
		Limit:      int32(limit),
	}

	transfers, err := a.transferStore.GetTransfersBetweenUsers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("no transfers was found between %d user and %d user", from_uid, to_uid), "no transfers found")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}
