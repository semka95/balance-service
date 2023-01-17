package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/semka95/balance-service/domain"
	transferModel "github.com/semka95/balance-service/transfer/repository"
	userModel "github.com/semka95/balance-service/user/repository"
)

// transferUcase represents usecase for transfer
type transferUcase struct {
	transferStore transferModel.Querier
	userStore     userModel.Querier
	db            *sql.DB
}

// New will create new transferUcase object representation of domain.TransferUsecase
func New(transferStore transferModel.Querier, userStore userModel.Querier, db *sql.DB) domain.TransferUsecase {
	return &transferUcase{
		transferStore: transferStore,
		userStore:     userStore,
		db:            db,
	}
}

// GetTransfer returns transfer by id
func (uc *transferUcase) GetTransfer(ctx context.Context, id int64) (*transferModel.Transfer, error) {
	transfer, err := uc.transferStore.GetTransferByID(ctx, id)
	return &transfer, err
}

// GetInboundTransfers returns all inbound transfers for user
func (uc *transferUcase) GetInboundTransfers(ctx context.Context, params transferModel.GetInboundTransfersParams) ([]transferModel.Transfer, error) {
	transfers, err := uc.transferStore.GetInboundTransfers(ctx, params)
	return transfers, err
}

// GetOutboundTransfers returns all outbound transfers for user
func (uc *transferUcase) GetOutboundTransfers(ctx context.Context, params transferModel.GetOutboundTransfersParams) ([]transferModel.Transfer, error) {
	transfers, err := uc.transferStore.GetOutboundTransfers(ctx, params)
	return transfers, err
}

// GetTransfersBetweenUsers returns all transfers between two users
func (uc *transferUcase) GetTransfersBetweenUsers(ctx context.Context, params transferModel.GetTransfersBetweenUsersParams) ([]transferModel.Transfer, error) {
	transfers, err := uc.transferStore.GetTransfersBetweenUsers(ctx, params)
	return transfers, err
}

// CreateTransfer creates transfer
func (uc *transferUcase) CreateTransfer(ctx context.Context, params transferModel.CreateTransferParams) (*transferModel.Transfer, error) {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback()

	userFrom, err := uc.userStore.GetUser(ctx, params.FromUserID)
	if err != nil {
		return nil, err
	}
	userTo, err := uc.userStore.GetUser(ctx, params.ToUserID)
	if err != nil {
		return nil, err
	}

	userFrom.Balance = userFrom.Balance.Sub(params.Amount)
	userTo.Balance = userTo.Balance.Add(params.Amount)

	_, err = uc.userStore.UpdateBalance(ctx, userModel.UpdateBalanceParams{ID: userFrom.ID, Balance: userFrom.Balance})
	if err != nil {
		return nil, err
	}

	_, err = uc.userStore.UpdateBalance(ctx, userModel.UpdateBalanceParams{ID: userTo.ID, Balance: userTo.Balance})
	if err != nil {
		return nil, err
	}

	transfer, err := uc.transferStore.CreateTransfer(ctx, params)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit transaction: %w", err)
	}

	return &transfer, nil
}
