package usecase

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/semka95/balance-service/domain"
	invoiceModel "github.com/semka95/balance-service/invoice/repository"
)

// invoiceUcase represents usecase for invoice
type invoiceUcase struct {
	invoiceStore invoiceModel.Querier
	db           *sql.DB
}

// New will create new an invoiceUcase object representation of domain.InvoiceUsecase interface
func New(invoiceStore invoiceModel.Querier, db *sql.DB) domain.InvoiceUsecase {
	return &invoiceUcase{
		invoiceStore: invoiceStore,
		db:           db,
	}
}

// GetInvoiceByID returns invoice by id
func (uc *invoiceUcase) GetInvoiceByID(ctx context.Context, id int64) (*invoiceModel.Invoice, error) {
	invoice, err := uc.invoiceStore.GetInvoiceByID(ctx, id)
	return &invoice, err
}

// GetInvoicesByUserID returns all invoices of user
func (uc *invoiceUcase) GetInvoicesByUserID(ctx context.Context, params invoiceModel.GetInvoicesByUserIDParams) ([]invoiceModel.Invoice, error) {
	invoice, err := uc.invoiceStore.GetInvoicesByUserID(ctx, params)
	return invoice, err
}

// CreateInvoice creates invoice
func (uc *invoiceUcase) CreateInvoice(ctx context.Context, params invoiceModel.CreateInvoiceParams) (*invoiceModel.Invoice, error) {
	invoice, err := uc.invoiceStore.CreateInvoice(ctx, params)
	return &invoice, err
}

// UpdateStatus updates invoice status
func (uc *invoiceUcase) UpdateStatus(ctx context.Context, params invoiceModel.UpdateStatusParams) (*invoiceModel.Invoice, error) {
	tx, err := uc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = uc.invoiceStore.GetInvoiceByID(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	invoice, err := uc.invoiceStore.UpdateStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("can't commit transaction: %w", err)
	}

	return &invoice, nil
}
