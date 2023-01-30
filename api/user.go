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

	"github.com/semka95/balance-service/domain"
	userModel "github.com/semka95/balance-service/user/repository"
)

// GET /user/{id} - returns user balance
func (a *API) getBalance(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}

	user, err := a.userUcase.GetUser(r.Context(), int64(userID))
	if err != nil {
		SendErrorJSON(w, r, domain.GetStatusCode(err), err, "can't get user balance")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, JSON{"balance": user.Balance.String()})
}

// PATCH /user/{id}/deposit - deposits money to user balance
func (a *API) depositMoney(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}

	params := userModel.UpdateBalanceParams{}
	if err = render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}
	params.ID = int64(userID)

	if params.Balance.IsNegative() || params.Balance.IsZero() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New("bad balance"), fmt.Sprintf("invalid balance: %s, should be greater than zero", params.Balance.String()))
		return
	}

	user, err := a.userUcase.UpdateBalance(r.Context(), params)
	if err != nil {
		SendErrorJSON(w, r, domain.GetStatusCode(err), err, "can't deposit money to user")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// PUT /user/{id}/withdraw - withdraws money from user balance
func (a *API) withdrawMoney(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid user id")
		return
	}

	params := userModel.UpdateBalanceParams{}
	if err = render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to balance")
		return
	}
	params.ID = int64(userID)

	if params.Balance.IsNegative() || params.Balance.IsZero() {
		SendErrorJSON(w, r, http.StatusBadRequest, errors.New(""), fmt.Sprintf("invalid balance: %s, should be greater then zero", params.Balance.String()))
		return
	}

	params.Balance = params.Balance.Neg()
	user, err := a.userUcase.UpdateBalance(r.Context(), params)
	if errors.Is(err, sql.ErrNoRows) {
		SendErrorJSON(w, r, http.StatusNotFound, err, "user not found")
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "")
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, user)
}

// POST /user - create user
func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	params := userModel.CreateUserParams{}

	if err := render.DecodeJSON(r.Body, &params); err != nil {
		SendErrorJSON(w, r, http.StatusBadRequest, err, "invalid request body, can't decode it to user")
		return
	}

	user, err := a.userUcase.CreateUser(r.Context(), params)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
		SendErrorJSON(w, r, http.StatusBadRequest, err, fmt.Sprintf("user with email '%s' already exists", params.Email))
		return
	}
	if err != nil {
		SendErrorJSON(w, r, http.StatusInternalServerError, err, "can't create user")
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, user)
}
