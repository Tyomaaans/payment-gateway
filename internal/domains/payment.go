package domains

import (
	"time"
)

type Status string

const (
	StatusPaid      Status = "paid"
	StatusPending   Status = "pending"
	StatusCancelled Status = "cancelled"
)

type MerchantEntity struct {
	ID           string
	DokuClientID string
	SecretKey    string
	IsSandbox    bool
}

type PaymentOrderEntity struct {
	ID              string
	InvoiceNumber   string
	MerchantID      string
	Amount          int64
	Currency        string
	PaymentMethod   *string
	Status          Status
	DokuTransID     *string
	ExpiredAt       *time.Time
	PaidAt          *time.Time
}

type PaymentTransactionEntity struct {
	ID              string
	PaymentOrderID  string
	DokuRequestID   string
	DokuResponseID  *string
	RequestPayload  string
	ResponsePayload *string
	HttpStatusCode  *int
}

type PaymentWebhookLogEntity struct {
	ID              string
	PaymentOrderID  *string
	Signature       string
	IsValidSig      bool
	RawPayload      string
	ProcessedStatus Status
	ReceivedAt      time.Time
}

type IdempotencyKeyEntity struct {
	Key         string    `gorm:"primaryKey;type:varchar(255)"`
	RequestHash string    `gorm:"type:varchar(255);not null"`
	ResponseBody string   `gorm:"type:jsonb"`
	CreatedAt   time.Time `gorm:"not null"`
	ExpiredAt   time.Time `gorm:"not null;index"`
}