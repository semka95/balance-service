// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package repository

import (
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type ValidStatus string

const (
	ValidStatusNew      ValidStatus = "new"
	ValidStatusAccepted ValidStatus = "accepted"
	ValidStatusRejected ValidStatus = "rejected"
	ValidStatusError    ValidStatus = "error"
)

func (e *ValidStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ValidStatus(s)
	case string:
		*e = ValidStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for ValidStatus: %T", src)
	}
	return nil
}

type NullValidStatus struct {
	ValidStatus ValidStatus
	Valid       bool // Valid is true if ValidStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullValidStatus) Scan(value interface{}) error {
	if value == nil {
		ns.ValidStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ValidStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullValidStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.ValidStatus, nil
}

type Invoice struct {
	ID            int64           `json:"id"`
	ServiceID     int64           `json:"service_id"`
	OrderID       int64           `json:"order_id"`
	UserID        int64           `json:"user_id"`
	Amount        decimal.Decimal `json:"amount"`
	PaymentStatus ValidStatus     `json:"payment_status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type Transfer struct {
	ID         int64     `json:"id"`
	FromUserID int64     `json:"from_user_id"`
	ToUserID   int64     `json:"to_user_id"`
	Amount     string    `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Balance   string    `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
