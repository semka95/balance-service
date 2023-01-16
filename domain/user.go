package domain

import (
	"context"

	userModel "github.com/semka95/balance-service/user/repository"
)

// UserUsecase represents User's usecases
type UserUsecase interface {
	CreateUser(ctx context.Context, params userModel.CreateUserParams) (*userModel.User, error)
	GetUser(ctx context.Context, id int64) (*userModel.User, error)
	UpdateBalance(ctx context.Context, params userModel.UpdateBalanceParams) (*userModel.User, error)
}
