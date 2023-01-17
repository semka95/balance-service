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
	userModel "github.com/semka95/balance-service/user/repository"
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

	transfers, err := a.transferStore.GetTransfersBetweenUsers(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get transfers")
		return
	}
	if len(transfers) == 0 {
		SendErrorJSON(w, r, http.StatusBadRequest, fmt.Errorf("no transfers was found between %d user and %d user", fromUID, toUID), "no transfers found")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfers)
}

// POST /transfer - transfers money from one user to another
func (a *API) transfer(w http.ResponseWriter, r *http.Request) {
	params := userModel.TransferMoneyParams{}
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}

	tx, err := a.db.BeginTx(r.Context(), nil)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't start transaction")
		return
	}
	defer tx.Rollback()

	userFrom, err := a.userStore.GetUser(r.Context(), params.FromUserID)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}
	userTo, err := a.userStore.GetUser(r.Context(), params.ToUserID)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	userFrom.Balance = userFrom.Balance.Sub(params.Amount)
	// TODO: maybe redundant, because database ensures it's not negative
	if userFrom.Balance.IsNegative() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New(""), fmt.Sprintf("not enough money on balance, available only %s", userFrom.Balance))
		return
	}
	userTo.Balance = userTo.Balance.Add(params.Amount)

	_, err = a.userStore.UpdateBalance(r.Context(), userModel.UpdateBalanceParams{ID: userFrom.ID, Balance: userFrom.Balance})
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	_, err = a.userStore.UpdateBalance(r.Context(), userModel.UpdateBalanceParams{ID: userTo.ID, Balance: userTo.Balance})
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	trParams := transferModel.CreateTransferParams{
		FromUserID: userFrom.ID,
		ToUserID:   userTo.ID,
		Amount:     params.Amount,
	}
	transfer, err := a.transferStore.CreateTransfer(r.Context(), trParams)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create transfer record")
		return
	}

	if err := tx.Commit(); err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't commit transaction")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, transfer)
}
