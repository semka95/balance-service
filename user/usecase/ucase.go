package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/semka95/balance-service/domain"
	userModel "github.com/semka95/balance-service/user/repository"
)

// userUcase represents usecase for user
type userUcase struct {
	userStore userModel.Querier
	db        *sql.DB
}

// New will create new an userUcase object representation of domain.UserUsecase interface
func New(userStore userModel.Querier, db *sql.DB) domain.UserUsecase {
	return &userUcase{
		userStore: userStore,
		db:        db,
	}
}

// GetUser returns User by id
func (uc *userUcase) GetUser(ctx context.Context, id int64) (*userModel.User, error) {
	user, err := uc.userStore.GetUser(ctx, id)

	return &user, err
}

// CreateUser creates new user
func (uc *userUcase) CreateUser(ctx context.Context, params userModel.CreateUserParams) (*userModel.User, error) {
	params.Balance = decimal.NewFromInt(0)

	user, err := uc.userStore.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateBalance updates user balance
func (uc *userUcase) UpdateBalance(ctx context.Context, params userModel.UpdateBalanceParams) (*userModel.User, error) {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback()

	user, err := uc.userStore.GetUser(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	params.Balance = params.Balance.Add(user.Balance)
	user, err = uc.userStore.UpdateBalance(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("can't update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit transaction: %w", err)
	}

	return &user, nil
}
