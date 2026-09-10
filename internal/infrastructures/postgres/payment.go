package postgres

import (
	"time"

	"payment-gateway/internal/domains"
)

type MerchantStorage struct {
	ID           string    `gorm:"primaryKey"`
	DokuClientID string    `gorm:"type:varchar(255);not null"`
	SecretKey    string    `gorm:"type:text;not null"` // encrypted at rest
	IsSandbox    bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

type PaymentOrderStorage struct {
	ID              string          `gorm:"primaryKey"`
	InvoiceNumber   string          `gorm:"type:varchar(100);not null;uniqueIndex"`
	MerchantID      string          `gorm:"type:varchar(255);not null;index"`
	Amount          int64           `gorm:"not null"` // simpan dalam smallest unit (sen)
	Currency        string          `gorm:"type:varchar(10);not null;default:'IDR'"`
	PaymentMethod   *string         `gorm:"type:varchar(50)"` // VA, QRIS, CC, dll
	Status          domains.Status  `gorm:"type:varchar(50);not null;index"`
	DokuTransID     *string         `gorm:"type:varchar(255);index"`
	ExpiredAt       *time.Time      `gorm:"index"`
	PaidAt          *time.Time      `gorm:"index"`
	CreatedAt       time.Time       `gorm:"not null"`
	UpdatedAt       time.Time       `gorm:"not null;index"`
}

type PaymentTransactionStorage struct {
	ID              string     `gorm:"primaryKey"`
	PaymentOrderID  string     `gorm:"type:varchar(255);not null;index"`
	DokuRequestID   string     `gorm:"type:varchar(255);not null"`
	DokuResponseID  *string    `gorm:"type:varchar(255)"`
	RequestPayload  string     `gorm:"type:jsonb;not null"`
	ResponsePayload *string    `gorm:"type:jsonb"`
	HttpStatusCode  *int       `gorm:"type:smallint"`
	CreatedAt       time.Time  `gorm:"not null;index"`
}

type PaymentWebhookLogStorage struct {
	ID              string         `gorm:"primaryKey"`
	PaymentOrderID  *string        `gorm:"type:varchar(255);index"`
	Signature       string         `gorm:"type:text;not null"`
	IsValidSig      bool           `gorm:"not null"`
	RawPayload      string         `gorm:"type:jsonb;not null"`
	ProcessedStatus domains.Status `gorm:"type:varchar(50);not null"`
	ReceivedAt      time.Time      `gorm:"not null;index"`
}