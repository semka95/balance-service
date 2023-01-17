package domain

import (
	"context"

	transferModel "github.com/semka95/balance-service/transfer/repository"
)

// TransferUsecase represents Transfer's usecases
type TransferUsecase interface {
	GetTransfer(ctx context.Context, id int64) (*transferModel.Transfer, error)
	GetInboundTransfers(ctx context.Context, params transferModel.GetInboundTransfersParams) ([]transferModel.Transfer, error)
	GetOutboundTransfers(ctx context.Context, params transferModel.GetOutboundTransfersParams) ([]transferModel.Transfer, error)
	GetTransfersBetweenUsers(ctx context.Context, params transferModel.GetTransfersBetweenUsersParams) ([]transferModel.Transfer, error)
	CreateTransfer(ctx context.Context, params transferModel.CreateTransferParams) (*transferModel.Transfer, error)
}
