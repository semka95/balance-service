package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/shopspring/decimal"

	userModel "github.com/semka95/balance-service/user/repository"
)

// GET /user/{id} - returns user balance
func (a *API) getBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}

	user, err := a.userStore.GetUser(r.Context(), int64(userID))
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, fmt.Sprintf("user with %d id not found", userID))
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't get balance")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, JSON{"balance": user.Balance.String()})
}

// PUT /user/deposit - deposits money to user balance
func (a *API) depositMoney(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}
	params := userModel.UpdateBalanceParams{}
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}
	params.ID = int64(userID)

	if params.Balance.IsNegative() || params.Balance.IsZero() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New(""), fmt.Sprintf("invalid balance: %s, should be greater then zero", params.Balance.String()))
		return
	}

	tx, err := a.db.BeginTx(r.Context(), nil)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't start transaction")
		return
	}
	defer tx.Rollback()
	user, err := a.userStore.GetUser(r.Context(), params.ID)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	params.Balance = params.Balance.Add(user.Balance)
	rows, err := a.userStore.UpdateBalance(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	if err := tx.Commit(); err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't commit transaction")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, rows)
}

// PUT /user/withdraw - withdraws money from user balance
func (a *API) withdrawMoney(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}
	params := userModel.UpdateBalanceParams{}
	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}
	params.ID = int64(userID)

	if params.Balance.IsNegative() || params.Balance.IsZero() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New(""), fmt.Sprintf("invalid balance: %s, should be greater then zero", params.Balance.String()))
		return
	}

	tx, err := a.db.BeginTx(r.Context(), nil)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't start transaction")
		return
	}
	defer tx.Rollback()
	user, err := a.userStore.GetUser(r.Context(), params.ID)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	params.Balance = user.Balance.Sub(params.Balance)
	// TODO: maybe redundant, because database ensures it's not negative
	if params.Balance.IsNegative() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New(""), fmt.Sprintf("not enough money on balance, available only %s", user.Balance))
		return
	}

	rows, err := a.userStore.UpdateBalance(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't update balance")
		return
	}

	if err := tx.Commit(); err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't commit transaction")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, rows)
}

// POST /user - create user
func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	//TODO: check email uniqeness
	createUser := userModel.CreateUserParams{}

	if err := render.DecodeJSON(r.Body, &createUser); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to user")
		return
	}

	createUser.Balance = decimal.NewFromInt(0)

	user, err := a.userStore.CreateUser(r.Context(), createUser)
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create user")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}
