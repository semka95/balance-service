package domain

import (
	"context"

	invoiceModel "github.com/semka95/balance-service/invoice/repository"
)

// InvoiceUsecase represents Invoice's usecases
type InvoiceUsecase interface {
	GetInvoiceByID(ctx context.Context, id int64) (*invoiceModel.Invoice, error)
	GetInvoicesByUserID(ctx context.Context, params invoiceModel.GetInvoicesByUserIDParams) ([]invoiceModel.Invoice, error)
	CreateInvoice(ctx context.Context, params invoiceModel.CreateInvoiceParams) (*invoiceModel.Invoice, error)
	UpdateStatus(ctx context.Context, params invoiceModel.UpdateStatusParams) (*invoiceModel.Invoice, error)
}
